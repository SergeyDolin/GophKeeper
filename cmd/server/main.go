package main

import (
	"gophkeeper/internal/handlers"
	"log"
	"net/http"
)

func main() {
	h := handlers.NewHandler()

	log.Println("server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", h.Router()))
}
