package main

import (
    "log"
    "net/http"
    "github.com/gorilla/mux"
    "time"
    "github.com/joho/godotenv"
    "github.com/whosthefunkyy/go-rest-api-example/handlers"
    "github.com/whosthefunkyy/go-rest-api-example/middleware"
    "github.com/whosthefunkyy/go-rest-api-example/db"
    "github.com/whosthefunkyy/go-rest-api-example/repository"
)

func main() {
    err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	db.ConnectGorm()
	db.AutoMigrate()

	userRepo := &repository.GormUserRepository{DB: db.DB} // ✅

	r := mux.NewRouter()
	r.Use(middleware.WithTimeoutMiddleware(5 * time.Second))

	h := &handlers.Handler{Repo: userRepo} // ✅

	api := r.PathPrefix("/api/v1").Subrouter()

	api.HandleFunc("/users", h.GetUsers).Methods("GET")
	api.HandleFunc("/users/{id}", h.GetUser).Methods("GET")
	api.HandleFunc("/users", h.CreateUser).Methods("POST")
	api.HandleFunc("/users/{id}", h.UpdateUser).Methods("PUT")
	api.HandleFunc("/users/{id}", h.DeleteUser).Methods("DELETE")
	//api.HandleFunc("/users/{id}", h.PatchUser).Methods("PATCH")

	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", api))
}
