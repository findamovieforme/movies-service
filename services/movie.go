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
	GetMovies(params map[string]string) ([]models.Movie, error)
	GetTrendingMovies(page int) ([]models.Movie, error)
	GetRecentMovies(page int) ([]models.Movie, error)
	GetMovieGenres() ([]models.Genre, error)
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

func (s *MovieService) GetMovies(params map[string]string) ([]models.Movie, error) {
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
	return s.GetMovies(params)
}

func (s *MovieService) GetRecentMovies(page int) ([]models.Movie, error) {
	params := make(map[string]string)
	params["sort_by"] = "release_date.desc"
	params["release_date.lte"] = time.Now().Format("2006-01-02")
	// Popularity must be decent
	params["page"] = strconv.Itoa(page)
	params["vote_count.gte"] = "100"

	return s.GetMovies(params)
}

func (s *MovieService) GetMovieGenres() ([]models.Genre, error) {
	genres, err := s.API.GetMovieGenres(nil)
	if err != nil {
		log.Fatal(err)
	}

	genreList := make([]models.Genre, 0)
	// Convert the genres to our model
	for _, genre := range genres.Genres {
		genre := models.Genre{
			ID:   genre.ID,
			Name: genre.Name,
		}
		genreList = append(genreList, genre)
	}
	return genreList, nil
}

// Function to fetch movie data for a list of IMDb IDs
func (s *MovieService) GetMoviesByIMDbIDs(imdbIDs []string) ([]models.Movie, error) {
	var movies []models.Movie
	// for _, imdbID := range imdbIDs {
	// 	url := fmt.Sprintf("https://api.themoviedb.org/3/find/%s?api_key=%s&external_source=imdb_id", imdbID, s.apiKey)

	// 	// Call the API for each IMDb ID
	// 	s.API.GetMovieGenres()

	// 	resp, err := http.Get(url)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	defer resp.Body.Close()

	// 	// Decode the API response
	// 	var result struct {
	// 		MovieResults []models.Movie `json:"movie_results"`
	// 	}
	// 	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
	// 		return nil, err
	// 	}

	// 	// Append each result
	// 	if len(result.MovieResults) > 0 {
	// 		movies = append(movies, result.MovieResults[0])
	// 	}
	// }
	return movies, nil
}
