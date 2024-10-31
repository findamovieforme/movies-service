package services

import (
	"log"
	"strconv"

	"github.com/movierecuh/movies-service/helpers"
	"github.com/movierecuh/movies-service/models"
	"github.com/ryanbradynd05/go-tmdb"
)

type MovieServiceInterface interface {
	GetMovies() ([]models.Movie, error)
}

type MovieService struct {
	API *tmdb.TMDb
}

func GetMovieService() *MovieService {
	apiKey, err := helpers.LoadEnv("TMDB_API_KEY")
	if err != nil {
		log.Fatal(err)
	}

	config := tmdb.Config{
		APIKey:   apiKey,
		Proxies:  nil,
		UseProxy: false,
	}
	tmdbAPI := tmdb.Init(config)
	return &MovieService{API: tmdbAPI}
}

func (s *MovieService) GetMovies() ([]models.Movie, error) {
	apiRes, err := s.API.DiscoverMovie(nil)
	if err != nil {
		return nil, err
	}

	var movies []models.Movie
	for _, movie := range apiRes.Results {
		movies = append(movies, models.Movie{
			ID:    movie.ID,
			Title: movie.Title,
			Year: func() int {
				year, err := strconv.Atoi(movie.ReleaseDate[:4])
				if err != nil {
					log.Println("Error converting year:", err)
					return 0
				}
				return year
			}(),
			Poster: movie.PosterPath,
		})
	}
	return movies, nil
}
