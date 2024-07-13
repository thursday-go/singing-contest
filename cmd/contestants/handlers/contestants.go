package handlers

import (
	"net/http"
	"encoding/json"

	"contestants-service/db"
)

func GetContestants(w http.ResponseWriter, r *http.Request) {
	contestants, err := db.FetchContestantsFromDB()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(contestants)
}
