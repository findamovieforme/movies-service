package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand/v2"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sagemakerruntime"
	"github.com/findamovieforme/movies-service/helpers"
	"github.com/findamovieforme/movies-service/models"
	"github.com/openai/openai-go" // imported as openai
	"github.com/openai/openai-go/option"
	"github.com/ryanbradynd05/go-tmdb"
)

type MovieServiceInterface interface {
	GetMovies(params map[string]string) ([]models.Movie, error)
	GetTrendingMovies(page int, genreID int) ([]models.Movie, error)
	GetRecentMovies(page int) ([]models.Movie, error)
	GetMovieGenres() ([]models.Genre, error)
	GetRecommendations(movieName string) ([]models.Movie, error)
	GetRecommendationsGrouped(movieNames []string) ([]models.Movie, error)

	GetMovieDetails(movieID int) (models.ExtendedMovie, error)
	GetMovieDetailsByTitle(movieTitle string) (models.Movie, error)
	SearchMovie(movieTitle string) ([]models.Movie, error)

	GetGptResponse(prompt string) ([]models.Movie, error)
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

var openaiAPIKey string

func GetMovieService() *MovieService {
	apiKey, err := helpers.LoadEnv("TMDB_API_KEY")
	if err != nil {
		log.Fatal(err)
	}
	openai, err := helpers.LoadEnv("OPENAI_API_KEY")
	if err != nil {
		log.Fatal(err)
	}
	openaiAPIKey = openai

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

func (s *MovieService) GetTrendingMovies(page int, genreID int) ([]models.Movie, error) {
	params := make(map[string]string)
	params["page"] = strconv.Itoa(page)
	params["sort_by"] = "popularity.desc"
	if genreID != 0 {
		params["with_genres"] = strconv.Itoa(genreID)
	}
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
	genreList := getGenresWithTrendingMovies(s, genres)

	return genreList, nil
}

func (s *MovieService) GetRecommendations(movieName string) ([]models.Movie, error) {

	// valkeyAddr := "valkeycluster-0001-001.valkeycluster.9eytty.use2.cache.amazonaws.com:6379" // Replace with your actual Valkey endpoint
	// valkeyClient, err := valkey.NewClient(valkey.ClientOption{
	// 	InitAddress: []string{valkeyAddr},
	// })
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// ctx := context.Background()
	// cacheKey := fmt.Sprintf("recommendations:%s", movieName)

	// // Attempt to retrieve cached recommendations
	// resp := valkeyClient.Do(ctx, valkeyClient.B().Get().Key(cacheKey).Build())
	// if err := resp.Error(); err == nil {
	// 	cachedData, _ := resp.ToString()
	// 	if cachedData != "" {
	// 		var movieRecommendations []models.Movie
	// 		if err := json.Unmarshal([]byte(cachedData), &movieRecommendations); err == nil {
	// 			log.Printf("Cache hit for movie: %s", movieName)
	// 			return movieRecommendations, nil
	// 		}
	// 		log.Printf("Error unmarshalling cached data: %v", err)
	// 	}
	// } else if !valkey.IsValkeyNil(err) {
	// 	log.Printf("Error retrieving from cache: %v", err)
	// }
	// Load the AWS configuration
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-2"))
	if err != nil {
		log.Fatalf("Unable to load SDK config, %v", err)
	}

	// Create a SageMaker Runtime client
	svc := sagemakerruntime.NewFromConfig(cfg)

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
			// Update the image URLs
			movie.BackdropPath = tmdbImageBaseURL + movie.BackdropPath
			movie.PosterPath = tmdbImageBaseURL + movie.PosterPath
			movieRecommendations = append(movieRecommendations, movie)
		}

		// recommendationsData, err := json.Marshal(movieRecommendations)
		// if err == nil {
		// 	setResp := valkeyClient.Do(ctx, valkeyClient.B().Set().Key(cacheKey).Value(string(recommendationsData)).ExSeconds(3600).Build())
		// 	if setResp.Error() != nil {
		// 		log.Printf("Error caching recommendations: %v", setResp.Error())
		// 	}
		// } else {
		// 	log.Printf("Error marshalling recommendations: %v", err)
		// }

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

func (s *MovieService) GetRecommendationsGrouped(movieNames []string) ([]models.Movie, error) {
	var movies []models.Movie
	for _, movieName := range movieNames {
		res, err := s.GetRecommendations(movieName)
		if err != nil {
			fmt.Println("Error getting recommendation for movie ", movieName, err)
			continue
		}
		movies = append(movies, res...)
	}
	// Shuffle the movies array
	for i := range movies {
		j := rand.IntN(i + 1)
		movies[i], movies[j] = movies[j], movies[i]
	}
	return movies, nil
}

func (s *MovieService) GetMovieDetails(movieID int) (models.ExtendedMovie, error) {
	movie, err := s.API.GetMovieInfo(movieID, nil)
	fmt.Println(movie)
	if err != nil {
		return models.ExtendedMovie{}, err
	}

	// Update the image URLs
	movie.BackdropPath = tmdbImageBaseURL + movie.BackdropPath
	movie.PosterPath = tmdbImageBaseURL + movie.PosterPath

	vids, err := s.API.GetMovieVideos(movieID, nil)
	if err != nil {
		fmt.Println("Error getting movie video: ", err)
	}

	// Get the trailer key
	var trailerKey string
	if vids != nil {
		for _, vid := range vids.Results {
			if vid.Type == "Trailer" && vid.Site == "YouTube" {
				trailerKey = vid.Key
			}
		}
	}
	movieTrimmed := models.ExtendedMovie{
		TrailerKey:    trailerKey,
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
		IMDBID:        movie.ImdbID,
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

func (s *MovieService) SearchMovie(movieTitle string) ([]models.Movie, error) {
	movies, err := s.API.SearchMovie(movieTitle, nil)
	if err != nil {
		return nil, err
	}

	// Update the image URLs
	for i := range movies.Results {
		movies.Results[i].BackdropPath = tmdbImageBaseURL + movies.Results[i].BackdropPath
		movies.Results[i].PosterPath = tmdbImageBaseURL + movies.Results[i].PosterPath
	}
	return movies.Results, nil
}

func getGenresWithTrendingMovies(s *MovieService, genres *tmdb.Genre) []models.Genre {
	var genreList []models.Genre
	// var mu sync.Mutex     // to ensure safe access to genreList
	// var wg sync.WaitGroup // to wait for all goroutines to complete

	for _, genre := range genres.Genres {
		// wg.Add(1)
		fmt.Println(genre.Name)
		// go func(genre struct {
		// 	ID   int
		// 	Name string
		// }) {
		// 	defer wg.Done()
		movie, err := s.GetTrendingMovies(1, genre.ID)

		// Prepare the genre struct without PosterPath in case of an error
		genreModel := models.Genre{
			ID:   genre.ID,
			Name: genre.Name,
		}

		// Only set PosterPath if there's no error and movies are available
		if err == nil && len(movie) > 0 {
			// Get random index
			randIndex := rand.IntN(len(movie))
			genreModel.PosterPath = movie[randIndex].PosterPath
		}
		// Append the genre model to the list safely
		// mu.Lock()
		genreList = append(genreList, genreModel)
		// mu.Unlock()
		// }(genre)
	}

	// Wait for all goroutines to complete
	// wg.Wait()
	return genreList
}

func (s *MovieService) GetGptResponse(userPrompt string) ([]models.Movie, error) {
	client := openai.NewClient(
		option.WithAPIKey(openaiAPIKey),
	)

	// Pre-built instruction to GPT
	preBuiltPrompt := fmt.Sprintf(`You are a movie recommendation assistant. The user will describe what they liked about a movie, focusing on themes, vibes, or specific aspects. Suggest up to 5 movies that match these described themes or elements. The user can also describe a story, you will have to find the movie and similar movies after that. Respond only with a JSON array of movie titles in this format:
[
  "Movie Title 1",
  "Movie Title 2",
  "Movie Title 3"
]

User's prompt: %s`, userPrompt)

	// Call GPT API with the pre-built instruction
	chatCompletion, err := client.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(preBuiltPrompt),
		}),
		Model:     openai.F(openai.ChatModelGPT4oMini), // GPT-4 or 3.5
		MaxTokens: openai.Int(200),                     // Allow space for a JSON array
	})
	fmt.Println(chatCompletion.Choices[0].Message.Content)
	if err != nil {
		return nil, err
	}

	// Parse GPT response (JSON array of movie titles)
	var movieTitles []string
	err = json.Unmarshal([]byte(chatCompletion.Choices[0].Message.Content), &movieTitles)
	if err != nil {
		return nil, fmt.Errorf("failed to parse GPT response: %v", err)
	}

	// Query TMDB for movie details
	var movies []tmdb.MovieShort
	for _, title := range movieTitles {
		tmdbResults, err := s.GetMovieDetailsByTitle(title)
		if err != nil {
			fmt.Println("Error getting movie details for ", title, err)
			continue // Skip movies that fail
		}

		// Update the image URLs
		tmdbResults.BackdropPath = tmdbImageBaseURL + tmdbResults.BackdropPath
		tmdbResults.PosterPath = tmdbImageBaseURL + tmdbResults.PosterPath
		movies = append(movies, tmdbResults)
	}

	return movies, nil
}
