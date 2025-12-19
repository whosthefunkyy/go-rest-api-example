package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/whosthefunkyy/go-rest-api-example/db"
	"github.com/whosthefunkyy/go-rest-api-example/handlers"
	"github.com/whosthefunkyy/go-rest-api-example/repository"
	"github.com/whosthefunkyy/go-rest-api-example/middleware" // если используешь таймауты
)

func main() {
	// 1. Подключение к БД
	db.ConnectGorm()
	db.AutoMigrate()

	// 2. Инициализация репозитория и хендлера
	userRepo := &repository.GormUserRepository{DB: db.DB}
	h := &handlers.Handler{Repo: userRepo}

	r := mux.NewRouter()

	// 3. Стандартные проверки
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "OK")
	}).Methods("GET")

	// 4. Группировка маршрутов API
	api := r.PathPrefix("/api/v1").Subrouter()
	
	// Добавляем мидлвар таймаута (опционально, но полезно)
	api.Use(middleware.WithTimeoutMiddleware(10 * time.Second))

	// РЕГИСТРАЦИЯ МАРШРУТОВ (то, чего у тебя не хватало)
	api.HandleFunc("/users", h.GetAllUsers).Methods("GET")
	api.HandleFunc("/users/{id}", h.GetUserByID).Methods("GET")
	api.HandleFunc("/users", h.CreateUser).Methods("POST")
	api.HandleFunc("/users/{id}", h.UpdateUser).Methods("PUT")
	api.HandleFunc("/users/{id}", h.DeleteUser).Methods("DELETE")

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "API is running WITH Database! Try /api/v1/users")
	}).Methods("GET")

	log.Println("Server started at :5000")
	if err := http.ListenAndServe(":5000", r); err != nil {
		log.Fatal(err)
	}
}
