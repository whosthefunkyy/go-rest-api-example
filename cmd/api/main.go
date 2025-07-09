package main

import (
    "log"
    "net/http"
    "github.com/gorilla/mux"
    "time"

    "github.com/whosthefunkyy/go-rest-api-example/handlers"
    "github.com/whosthefunkyy/go-rest-api-example/middleware"
    "github.com/whosthefunkyy/go-rest-api-example/db"
)

func main() {
    db.InitDB() // db.go
    defer db.DB.Close()
    r := mux.NewRouter()
    r.Use(middleware.WithTimeoutMiddleware(5 * time.Second)) 
    api := r.PathPrefix("/api/v1").Subrouter()
    h := &handlers.Handler{DB: db.DB}

    api.HandleFunc("/users", h.GetUsers).Methods("GET")
    api.HandleFunc("/users/{id}",handlers.GetUser).Methods("GET")
    api.HandleFunc("/users", handlers.CreateUser).Methods("POST")
    api.HandleFunc("/users/{id}", handlers.UpdateUser).Methods("PUT")
    api.HandleFunc("/users/{id}",handlers.DeleteUser).Methods("DELETE")
    api.HandleFunc("/users/{id}", handlers.PatchUser).Methods("PATCH")

    log.Println("Server started at :8080")
    log.Fatal(http.ListenAndServe(":8080", api))
}
