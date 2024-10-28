package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/movierecuh/movies-service/services"
)

type RecommenderHandler struct {
	Service services.RecommenderServiceInterface
}

func NewRecommenderHandler(service services.RecommenderServiceInterface) *RecommenderHandler {
	return &RecommenderHandler{Service: service}
}

func (h *RecommenderHandler) FetchRecommendations(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("userID")
	if userID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	recommendations, err := h.Service.GetRecommendations(userID)
	if err != nil {
		http.Error(w, "Failed to fetch recommendations", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(recommendations)
}
