package services

import (
	"log"
	"strconv"
	"time"

	"github.com/movierecuh/movies-service/helpers"
	"github.com/movierecuh/movies-service/models"
	"github.com/ryanbradynd05/go-tmdb"
)

type MovieServiceInterface interface {
	GetMovies(page int, params map[string]string) ([]models.Movie, error)
	GetTrendingMovies(page int) ([]models.Movie, error)
	GetRecentMovies(page int) ([]models.Movie, error)
}

type MovieService struct {
	API *tmdb.TMDb
}

const tmdbImageBaseURL = "https://image.tmdb.org/t/p/w500"

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

func (s *MovieService) GetMovies(page int, params map[string]string) ([]models.Movie, error) {
	// Call the API with the parameters
	apiRes, err := s.API.DiscoverMovie(params)
	if err != nil {
		return nil, err
	}

	// Update the image URLs
	for i := range apiRes.Results {
		apiRes.Results[i].BackdropPath = tmdbImageBaseURL + apiRes.Results[i].BackdropPath
		apiRes.Results[i].PosterPath = tmdbImageBaseURL + apiRes.Results[i].PosterPath
	}

	return apiRes.Results, nil
}

func (s *MovieService) GetTrendingMovies(page int) ([]models.Movie, error) {
	params := make(map[string]string)
	params["page"] = strconv.Itoa(page)
	params["sort_by"] = "popularity.desc"
	return s.GetMovies(page, params)
}

func (s *MovieService) GetRecentMovies(page int) ([]models.Movie, error) {
	params := make(map[string]string)
	params["sort_by"] = "release_date.desc"
	params["release_date.lte"] = time.Now().Format("2006-01-02")
	// Popularity must be decent
	params["vote_count.gte"] = "100"

	return s.GetMovies(page, params)
}
