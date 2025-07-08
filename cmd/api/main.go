package main

import (
    "log"
    "net/http"
    "github.com/gorilla/mux"
    "time"

    "myapi/handlers"
    "myapi/middleware"
    "myapi/db"
)

func main() {
    db.InitDB() // из db.go

    r := mux.NewRouter()
    r.Use(middleware.WithTimeoutMiddleware(5 * time.Second)) 
    api := r.PathPrefix("/api/v1").Subrouter()

    api.HandleFunc("/users", handlers.GetUsers).Methods("GET")
    api.HandleFunc("/users/{id}",handlers.GetUser).Methods("GET")
    api.HandleFunc("/users", handlers.CreateUser).Methods("POST")
    api.HandleFunc("/users/{id}", handlers.UpdateUser).Methods("PUT")
    api.HandleFunc("/users/{id}",handlers.DeleteUser).Methods("DELETE")
    api.HandleFunc("/users/{id}", handlers.PatchUser).Methods("PATCH")

    log.Println("Server started at :8080")
    log.Fatal(http.ListenAndServe(":8080", api))
}
