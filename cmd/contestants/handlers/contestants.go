package handlers

import (
    "net/http"
    "encoding/json"

    "contestants-service/models"
)

var contestants []models.Contestant

func init() {
    // Initialize with some dummy data
    contestants = []models.Contestant{
        {ID: "1", Name: "Contestant A", Location: "Location A"},
        {ID: "2", Name: "Contestant B", Location: "Location B"},
        {ID: "3", Name: "Contestant C", Location: "Location C"},
    }
}

func GetContestants(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(contestants)
}
