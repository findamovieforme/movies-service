// All handlers to get movies

package handlers

import (
	"net/http"

	"github.com/movierecuh/movies-service/models"
	"github.com/movierecuh/movies-service/services"
)

type MovieHandler struct {
	Service services.MovieServiceInterface
}

func NewMovieHandler(service services.MovieServiceInterface) *MovieHandler {
	return &MovieHandler{Service: service}
}

func (h *MovieHandler) Fetchmovies() AppHandler {
	return func(w http.ResponseWriter, r *http.Request) ([]models.Movie, error) {
		return h.Service.GetMovies()
	}
}
