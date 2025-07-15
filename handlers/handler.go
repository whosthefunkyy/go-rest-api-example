package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	//"strings"
	"log"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"context"
  
	"github.com/whosthefunkyy/go-rest-api-example/models"
	"github.com/whosthefunkyy/go-rest-api-example/utils"
	"github.com/whosthefunkyy/go-rest-api-example/hateoas"
	"github.com/whosthefunkyy/go-rest-api-example/repository"

)
type Handler struct {
	Repo repository.UserRepository
}
	// GET ALL USERS
func (h *Handler) GetUsers(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()

    // Проверка отмены контекста (например, если клиент закрыл соединение)
    select {
    case <-ctx.Done():
        utils.SendError(w, "Request canceled", http.StatusRequestTimeout)
        return
    default:
    }

    // Получение пользователей из репозитория
    users, err := h.Repo.GetAll()
    if err != nil {
        log.Printf("GetUsers GetAll error: %v", err)
        utils.SendError(w, "Database error", http.StatusInternalServerError)
        return
    }

    // Формирование HATEOAS-ответа
    result := make([]map[string]interface{}, 0, len(users))
    for _, u := range users {
        result = append(result, hateoas.CreateUserResponse(u))
    }

    // Успешный ответ
    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(result); err != nil {
        log.Printf("GetUsers response encoding error: %v", err)
        utils.SendError(w, "Failed to encode response", http.StatusInternalServerError)
        return
    }
}
	// GET One USER by ID
func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()

    // Проверка отмены запроса
    select {
    case <-ctx.Done():
        utils.SendError(w, "Request canceled", http.StatusRequestTimeout)
        return
    default:
    }

    // Валидация ID
    idStr := mux.Vars(r)["id"]
    id, err := strconv.Atoi(idStr)
    if err != nil {
        utils.SendError(w, "Invalid ID", http.StatusBadRequest)
        return
    }

    // Получение пользователя
      user, err := h.Repo.GetByID(id)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            utils.SendError(w, "User not found", http.StatusNotFound)
        } else if ctx.Err() != nil {
            utils.SendError(w, "Request timeout", http.StatusGatewayTimeout)
        } else {
            utils.SendError(w, "Database error", http.StatusInternalServerError)
        }
        return
    }

    // Успешный ответ
    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(hateoas.CreateUserResponse(*user)); err != nil {
        log.Printf("GetUser response encoding error: %v", err)
        utils.SendError(w, "Failed to encode response", http.StatusInternalServerError)
    }
}

	// CREATE A NEW USER
func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()

    // Проверка отмены контекста
    select {
    case <-ctx.Done():
        utils.SendError(w, "Request canceled", http.StatusRequestTimeout)
        return
    default:
    }

    // Декодирование тела запроса
    var user models.User
    if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
        utils.SendError(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    // Валидация данных
    if user.Name == "" {
        utils.SendError(w, "Name cannot be empty", http.StatusBadRequest)
        return
    }

    // Создание пользователя
    if err := h.Repo.Create(&user); err != nil {
        log.Printf("CreateUser Create error: %v", err)

        // Пример: можно сюда вставить кастомную проверку уникальности email и т.п.
        if errors.Is(err, context.DeadlineExceeded) || errors.Is(ctx.Err(), context.DeadlineExceeded) {
            utils.SendError(w, "Request timeout", http.StatusGatewayTimeout)
            return
        }

        utils.SendError(w, "Database error", http.StatusInternalServerError)
        return
    }

    // Успешный ответ
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    if err := json.NewEncoder(w).Encode(hateoas.CreateUserResponse(user)); err != nil {
        log.Printf("CreateUser response encoding error: %v", err)
        utils.SendError(w, "Failed to encode response", http.StatusInternalServerError)
    }
}


	// UPDATE INFO into USER
func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()

    // Проверка отмены контекста
    select {
    case <-ctx.Done():
        utils.SendError(w, "Request canceled", http.StatusRequestTimeout)
        return
    default:
    }

    // Валидация ID
    idStr := mux.Vars(r)["id"]
    id, err := strconv.Atoi(idStr)
    if err != nil {
        utils.SendError(w, "Invalid ID", http.StatusBadRequest)
        return
    }

    // Декодирование тела запроса
    var user models.User
    if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
        utils.SendError(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    // Валидация данных
    if user.Name == "" {
        utils.SendError(w, "Name cannot be empty", http.StatusBadRequest)
        return
    }

    // Проверка существования пользователя
    existingUser, err := h.Repo.GetByID(id)
    if err != nil {
        switch {
        case errors.Is(err, gorm.ErrRecordNotFound):
            utils.SendError(w, "User not found", http.StatusNotFound)
        case errors.Is(err, context.DeadlineExceeded), errors.Is(ctx.Err(), context.DeadlineExceeded):
            utils.SendError(w, "Request timeout", http.StatusGatewayTimeout)
        default:
            log.Printf("UpdateUser GetByID error: %v", err)
            utils.SendError(w, "Database error", http.StatusInternalServerError)
        }
        return
    }

    // Обновление
    user.ID = existingUser.ID
    if err := h.Repo.Update(&user); err != nil {
        log.Printf("UpdateUser Update error: %v", err)
        if errors.Is(err, context.DeadlineExceeded) || errors.Is(ctx.Err(), context.DeadlineExceeded) {
            utils.SendError(w, "Request timeout", http.StatusGatewayTimeout)
            return
        }
        utils.SendError(w, "Database error", http.StatusInternalServerError)
        return
    }

    // Ответ
    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(hateoas.CreateUserResponse(user)); err != nil {
        log.Printf("UpdateUser response encoding error: %v", err)
        utils.SendError(w, "Failed to encode response", http.StatusInternalServerError)
    }
}

// DELETE USER
func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()

    // Проверка отмены контекста
    select {
    case <-ctx.Done():
        utils.SendError(w, "Request canceled", http.StatusRequestTimeout)
        return
    default:
    }

    // Валидация ID
    idStr := mux.Vars(r)["id"]
    id, err := strconv.Atoi(idStr)
    if err != nil {
        utils.SendError(w, "Invalid ID", http.StatusBadRequest)
        return
    }

    // Проверка существования пользователя
    user, err := h.Repo.GetByID(id)
    if err != nil {
        switch {
        case errors.Is(err, gorm.ErrRecordNotFound):
            utils.SendError(w, "User not found", http.StatusNotFound)
        case errors.Is(err, context.DeadlineExceeded), errors.Is(ctx.Err(), context.DeadlineExceeded):
            utils.SendError(w, "Request timeout", http.StatusGatewayTimeout)
        default:
            log.Printf("DeleteUser GetByID error: %v", err)
            utils.SendError(w, "Database error", http.StatusInternalServerError)
        }
        return
    }

    // Удаление
    if err := h.Repo.Delete(int(user.ID)); err != nil {
        log.Printf("DeleteUser Delete error: %v", err)
        if errors.Is(err, context.DeadlineExceeded) || errors.Is(ctx.Err(), context.DeadlineExceeded) {
            utils.SendError(w, "Request timeout", http.StatusGatewayTimeout)
            return
        }
        utils.SendError(w, "Database error", http.StatusInternalServerError)
        return
    }

    // Успешный ответ без тела
    w.WriteHeader(http.StatusNoContent)
}

