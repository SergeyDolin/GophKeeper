package main

import (
	"log"
	"main/internal/handlers"
	"net/http"
)

func main() {
	h := handlers.NewHandler()

	log.Println("server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", h.Router()))
}
