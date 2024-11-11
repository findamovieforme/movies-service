package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sagemakerruntime"
	"github.com/movierecuh/movies-service/helpers"
	"github.com/movierecuh/movies-service/models"
	"github.com/ryanbradynd05/go-tmdb"
)

type MovieServiceInterface interface {
	GetMovies(params map[string]string) ([]models.Movie, error)
	GetTrendingMovies(page int) ([]models.Movie, error)
	GetRecentMovies(page int) ([]models.Movie, error)
	GetMovieGenres() ([]models.Genre, error)
	GetRecommendations(movieName string) ([]models.Movie, error)
	GetMovieDetails(movieID int) (models.Movie, error)
	GetMovieDetailsByTitle(movieTitle string) (models.Movie, error)
}

func convertGenresToIDs(genres []struct {
	ID   int
	Name string
}) []int32 {
	var ids []int32
	for _, genre := range genres {
		ids = append(ids, int32(genre.ID))
	}
	return ids
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

func (s *MovieService) GetRecommendations(movieName string) ([]models.Movie, error) {

	// Load the AWS configuration
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-2"))
	if err != nil {
		log.Fatalf("Unable to load SDK config, %v", err)
	}

	// Create a SageMaker Runtime client
	svc := sagemakerruntime.NewFromConfig(cfg)

	// Specify the endpoint name and input payload
	endpointName := "movie-endpoint"
	inputPayload := fmt.Sprintf("{\"title\": \"%s\"}", movieName)

	// Call the SageMaker endpoint
	output, err := svc.InvokeEndpoint(context.TODO(), &sagemakerruntime.InvokeEndpointInput{
		EndpointName: aws.String(endpointName),
		ContentType:  aws.String("application/json"),
		Body:         []byte(inputPayload),
	})
	if err != nil {
		log.Printf("Error calling SageMaker endpoint: %v", err)
		return nil, err
	}

	// Parse the JSON response
	responseBody := output.Body
	var recommendationsResponse models.RecommendationsResponse
	var errorResponse models.ErrorResponse

	// First try to unmarshal into the recommendations response
	if err := json.Unmarshal(responseBody, &recommendationsResponse); err == nil && len(recommendationsResponse.Recommendations) > 0 {
		// Successfully parsed recommendations
		var movieRecommendations []models.Movie
		for _, title := range recommendationsResponse.Recommendations {
			movie, err := s.GetMovieDetailsByTitle(title)
			if err != nil {
				log.Printf("Error when getting movie details: %v", err)
				continue
			}
			movieRecommendations = append(movieRecommendations, movie)
		}
		return movieRecommendations, nil
	}

	// If no recommendations, try to unmarshal into the error response
	if err := json.Unmarshal(responseBody, &errorResponse); err == nil && errorResponse.Error != "" {
		log.Printf("Error when calling sagemaker: %v", err)
		return []models.Movie{}, nil
	}

	// If neither works, return a generic error
	log.Printf("Unexpected response when calling SageMaker endpoint: %v", err)
	return nil, fmt.Errorf("unexpected response from SageMaker: %s", string(responseBody))
}

func (s *MovieService) GetMovieDetails(movieID int) (models.Movie, error) {
	movie, err := s.API.GetMovieInfo(movieID, nil)
	fmt.Println(movie)
	if err != nil {
		return models.Movie{}, err
	}

	// Update the image URLs
	movie.BackdropPath = tmdbImageBaseURL + movie.BackdropPath
	movie.PosterPath = tmdbImageBaseURL + movie.PosterPath

	movieTrimmed := models.Movie{
		GenreIDs:      convertGenresToIDs(movie.Genres),
		Overview:      movie.Overview,
		ReleaseDate:   movie.ReleaseDate,
		BackdropPath:  movie.BackdropPath,
		PosterPath:    movie.PosterPath,
		Adult:         movie.Adult,
		ID:            movie.ID,
		OriginalTitle: movie.OriginalTitle,
		Popularity:    movie.Popularity,
		Video:         movie.Video,
		VoteCount:     movie.VoteCount,
		VoteAverage:   movie.VoteAverage,
		Title:         movie.Title,
	}

	return movieTrimmed, nil
}

func (s *MovieService) GetMovieDetailsByTitle(movieTitle string) (models.Movie, error) {
	movies, err := s.API.SearchMovie(movieTitle, nil)
	if err != nil {
		return models.Movie{}, err
	}
	movie := movies.Results[0]
	return movie, nil
}
