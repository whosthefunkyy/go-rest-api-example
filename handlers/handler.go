package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	//"strings"
	"context"
	"log"

	"github.com/gorilla/mux"
	"gorm.io/gorm"

	"github.com/whosthefunkyy/go-rest-api-example/hateoas"
	"github.com/whosthefunkyy/go-rest-api-example/repository"
	"github.com/whosthefunkyy/go-rest-api-example/utils"
)

type Handler struct {
	Repo repository.UserRepository
}
// GET ALL USERS
func (h *Handler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// ctx cancel
	select {
	case <-ctx.Done():
		utils.SendError(w, "Request canceled", http.StatusRequestTimeout)
		return
	default:
	}
	// get repo
	users, err := h.Repo.GetAll()
	if err != nil {
		log.Printf("GetUsers GetAll error: %v", err)
		utils.SendError(w, "Database error", http.StatusInternalServerError)
		return
	}
	// HATEOAS
	result := make([]map[string]interface{}, 0, len(users))
	for _, u := range users {
		result = append(result, hateoas.CreateUserResponse(u))
	}
	// success
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		log.Printf("GetUsers response encoding error: %v", err)
		utils.SendError(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
// GET One USER by ID
func (h *Handler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	select {
	case <-ctx.Done():
		utils.SendError(w, "Request canceled", http.StatusRequestTimeout)
		return
	default:
	}
	// valid ID
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.SendError(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	// get repo by id
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
	// success
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(hateoas.CreateUserResponse(*user)); err != nil {
		log.Printf("GetUser response encoding error: %v", err)
		utils.SendError(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// CREATE A NEW USER
func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// ctx cancel
	select {
	case <-ctx.Done():
		utils.SendError(w, "Request canceled", http.StatusRequestTimeout)
		return
	default:
	}
	// valid and pars
	user, err := utils.ParseAndValidateUser(r)
	if err != nil {
		switch err {
		case utils.ErrInvalidBody:
			utils.SendError(w, "Invalid request body", http.StatusBadRequest)
		case utils.ErrEmptyName:
			utils.SendError(w, "Name cannot be empty", http.StatusBadRequest)
		default:
			utils.SendError(w, "Validation error", http.StatusBadRequest)
		}
		return
	}
	// save into bd
	if err := h.Repo.Create(user); err != nil {
		log.Printf("CreateUser error: %v", err)
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(ctx.Err(), context.DeadlineExceeded) {
			utils.SendError(w, "Request timeout", http.StatusGatewayTimeout)
			return
		}
		utils.SendError(w, "Database error", http.StatusInternalServerError)
		return
	}
	// succes
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(hateoas.CreateUserResponse(*user)); err != nil {
		log.Printf("CreateUser response encoding error: %v", err)
		utils.SendError(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
// UPDATE INFO into USER
func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	//  ctx cancel
	select {
	case <-ctx.Done():
		utils.SendError(w, "Request canceled", http.StatusRequestTimeout)
		return
	default:
	}
	// valid and pars ID
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.SendError(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	// valid and pars body
	user, err := utils.ParseAndValidateUser(r)
	if err != nil {
		switch err {
		case utils.ErrInvalidBody:
			utils.SendError(w, "Invalid request body", http.StatusBadRequest)
		case utils.ErrEmptyName:
			utils.SendError(w, "Name cannot be empty", http.StatusBadRequest)
		default:
			utils.SendError(w, "Validation error", http.StatusBadRequest)
		}
		return 
	}
	// check existing user
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
	// update
	user.ID = existingUser.ID
	if err := h.Repo.Update(user); err != nil {
		log.Printf("UpdateUser Update error: %v", err)
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(ctx.Err(), context.DeadlineExceeded) {
			utils.SendError(w, "Request timeout", http.StatusGatewayTimeout)
			return
		}
		utils.SendError(w, "Database error", http.StatusInternalServerError)
		return
	}
	// answer
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(hateoas.CreateUserResponse(*user)); err != nil {
		log.Printf("UpdateUser response encoding error: %v", err)
		utils.SendError(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// DELETE USER
func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// ctx cancel
	select {
	case <-ctx.Done():
		utils.SendError(w, "Request canceled", http.StatusRequestTimeout)
		return
	default:
	}
	// valid ID
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.SendError(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	// chech existing user
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
	// delete
	if err := h.Repo.Delete(int(user.ID)); err != nil {
		log.Printf("DeleteUser Delete error: %v", err)
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(ctx.Err(), context.DeadlineExceeded) {
			utils.SendError(w, "Request timeout", http.StatusGatewayTimeout)
			return
		}
		utils.SendError(w, "Database error", http.StatusInternalServerError)
		return
	}
	// success
	w.WriteHeader(http.StatusNoContent)
}
