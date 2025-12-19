package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/whosthefunkyy/go-rest-api-example/db" 
)

func main() {

	db.ConnectGorm()
	db.AutoMigrate()

	r := mux.NewRouter()

	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "OK")
	}).Methods("GET")

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	
		fmt.Fprintf(w, "API is running WITH Database!")
	}).Methods("GET")

	log.Println("Server started at :5000")
	if err := http.ListenAndServe(":5000", r); err != nil {
		log.Fatal(err)
	}
}
