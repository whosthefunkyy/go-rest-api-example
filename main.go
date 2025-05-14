package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	// Define a handler function for our web server
	handler := func(w http.ResponseWriter, r *http.Request) {
		// Set the content type to plain text
		w.Header().Set("Content-Type", "text/plain")
		// Write "Hello, World!" as the response body
		fmt.Fprintln(w, "Hello, Stepik!")
	}

	// Register the handler function to handle all requests to the root ("/") path
	http.HandleFunc("/", handler)

	// Start the web server on port 8080
	log.Println("Server is running on :8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Println("Error:", err)
	}
}