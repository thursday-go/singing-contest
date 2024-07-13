package main

import (
	"net/http"
	"log"

	"contestants-service/handlers"
	"contestants-service/db"
)

func main() {
	db.InitDB()
	defer db.CloseDB()

	http.HandleFunc("/contestants", handlers.GetContestants)

	log.Println("Starting server on :8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
