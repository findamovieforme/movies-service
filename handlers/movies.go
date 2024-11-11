// All handlers to get movies

package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/movierecuh/movies-service/services"
)

type MovieHandler struct {
	Service services.MovieServiceInterface
}

func NewMovieHandler(service services.MovieServiceInterface) *MovieHandler {
	return &MovieHandler{Service: service}
}

func (h *MovieHandler) FetchTrendingMovies() AppHandler {
	return func(w http.ResponseWriter, r *http.Request) (interface{}, error) {
		pageInt := getPageFromRequest(r)
		return h.Service.GetTrendingMovies(pageInt)
	}
}

func (h *MovieHandler) FetchRecentlyReleasedMovies() AppHandler {
	return func(w http.ResponseWriter, r *http.Request) (interface{}, error) {
		pageInt := getPageFromRequest(r)
		return h.Service.GetRecentMovies(pageInt)
	}
}

func (h *MovieHandler) FetchGenres() AppHandler {
	return func(w http.ResponseWriter, r *http.Request) (interface{}, error) {
		return h.Service.GetMovieGenres()
	}
}

func (h *MovieHandler) FetchRecommendations() AppHandler {
	return func(w http.ResponseWriter, r *http.Request) (interface{}, error) {
		var req struct {
			Title string `json:"title"`
		}

		// Decode the JSON request body into the struct
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			return nil, err
		}

		// Use the extracted movie name to get recommendations
		return h.Service.GetRecommendations(req.Title)
	}
}

func (h *MovieHandler) FetchMovieDetails() AppHandler {
	return func(w http.ResponseWriter, r *http.Request) (interface{}, error) {

		// Extract movie name from query parameters as integer
		movieID, err := strconv.Atoi(r.URL.Query().Get("movieID"))
		if err != nil {
			return nil, err
		}
		// Use the extracted movie name to get recommendations
		return h.Service.GetMovieDetails(movieID)
	}
}

// Helper function to parse page number from request
func getPageFromRequest(r *http.Request) int {
	page := r.URL.Query().Get("page")
	pageInt := 1 // Default page
	if page != "" {
		if parsedPage, err := strconv.Atoi(page); err == nil {
			pageInt = parsedPage
		}
	}
	return pageInt
}
