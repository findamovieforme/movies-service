package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand/v2"
	"strconv"
	"time"

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
	log.Printf("[GetRecommendations] start movieName=%q", movieName)
	helpers.InitEnv()
	responseBody, err := helpers.CallLocalModel(movieName)
	if err != nil {
		log.Printf("[GetRecommendations] CallLocalModel failed for %q: %v", movieName, err)
		return nil, err
	}

	log.Printf("[GetRecommendations] got %d recommendation titles, resolving via TMDB", len(responseBody.Recommendations))

	// Successfully parsed recommendations
	var movieRecommendations []models.Movie

	for i, title := range responseBody.Recommendations {
		movie, err := s.GetMovieDetailsByTitle(title)
		if err != nil {
			log.Printf("[GetRecommendations] error getting movie details for recommendation %d title=%q: %v", i+1, title, err)
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

	log.Printf("[GetRecommendations] done movieName=%q: returning %d movies", movieName, len(movieRecommendations))
	return movieRecommendations, nil

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
