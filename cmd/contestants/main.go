package main

import (
    "net/http"
    "log"

    "contestants-service/handlers"
)

func main() {
    http.HandleFunc("/contestants", handlers.GetContestants)

    log.Println("Starting server on :8001")
    log.Fatal(http.ListenAndServe(":8001", nil))
}
