package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	// Наш главный эндпоинт для проверки здоровья
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "OK")
	}).Methods("GET")

	// Простой тестовый маршрут
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "API is running without DB for testing")
	}).Methods("GET")

	log.Println("Server started at :5000")
	// Важно: порт 8080 обязателен для EB Go платформы по умолчанию
	if err := http.ListenAndServe(":5000", r); err != nil {
		log.Fatal(err)
	}
}
