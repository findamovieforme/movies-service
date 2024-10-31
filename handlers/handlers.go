package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/movierecuh/movies-service/models"
)

type AppHandler func(http.ResponseWriter, *http.Request) ([]models.Movie, error)

func (fn AppHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	data, err := fn(w, r)
	if err != nil {
		http.Error(w, "Failed to process request", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
