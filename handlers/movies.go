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
		genreID := getIntFromRequest(r, "genreID")
		return h.Service.GetTrendingMovies(pageInt, genreID)
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
		// return h.Service.GetMovieGenres()
		genres := []map[string]interface{}{
			{
				"id":         28,
				"name":       "Action",
				"posterPath": "https://image.tmdb.org/t/p/w500/2uNW4WbgBXL25BAbXGLnLqX71Sw.jpg",
			},
			{
				"id":         12,
				"name":       "Adventure",
				"posterPath": "https://image.tmdb.org/t/p/w500/1g0dhYtq4irTY1GPXvft6k4YLjm.jpg",
			},
			{
				"id":         16,
				"name":       "Animation",
				"posterPath": "https://image.tmdb.org/t/p/w500/i77OInTKcrnRlAozFOaB6D5mk15.jpg",
			},
			{
				"id":         35,
				"name":       "Comedy",
				"posterPath": "https://image.tmdb.org/t/p/w500/cdqLnri3NEGcmfnqwk2TSIYtddg.jpg",
			},
			{
				"id":         80,
				"name":       "Crime",
				"posterPath": "https://image.tmdb.org/t/p/w500/j8Jx4vpBG258jPxd4o0kOLDwTvm.jpg",
			},
			{
				"id":         99,
				"name":       "Documentary",
				"posterPath": "https://image.tmdb.org/t/p/w500/bDQ95W5LPHW9FHlPj3QX3jvM9Z7.jpg",
			},
			{
				"id":         18,
				"name":       "Drama",
				"posterPath": "https://image.tmdb.org/t/p/w500/s2KE27j9cWDSJKYOg4KxXDwnoiv.jpg",
			},
			{
				"id":         10751,
				"name":       "Family",
				"posterPath": "https://image.tmdb.org/t/p/w500/yh64qw9mgXBvlaWDi7Q9tpUBAvH.jpg",
			},
			{
				"id":         14,
				"name":       "Fantasy",
				"posterPath": "https://image.tmdb.org/t/p/w500/dzDMewC0Hwv01SROiWgKOi4iOc1.jpg",
			},
			{
				"id":         36,
				"name":       "History",
				"posterPath": "https://image.tmdb.org/t/p/w500/dB6Krk806zeqd0YNp2ngQ9zXteH.jpg",
			},
			{
				"id":         27,
				"name":       "Horror",
				"posterPath": "https://image.tmdb.org/t/p/w500/l1175hgL5DoXnqeZQCcU3eZIdhX.jpg",
			},
			{
				"id":         10402,
				"name":       "Music",
				"posterPath": "https://image.tmdb.org/t/p/w500/x6wH1kowr6uNFJ12CVKsRHzC0cm.jpg",
			},
			{
				"id":         9648,
				"name":       "Mystery",
				"posterPath": "https://image.tmdb.org/t/p/w500/ycoXsJomPmPjtCfNweM0UWiTkPY.jpg",
			},
			{
				"id":         10749,
				"name":       "Romance",
				"posterPath": "https://image.tmdb.org/t/p/w500/f6PfAXtFEkJRcBtOjbzOgz8qqSK.jpg",
			},
			{
				"id":         878,
				"name":       "Science Fiction",
				"posterPath": "https://image.tmdb.org/t/p/w500/lqoMzCcZYEFK729d6qzt349fB4o.jpg",
			},
			{
				"id":         10770,
				"name":       "TV Movie",
				"posterPath": "https://image.tmdb.org/t/p/w500/65DkgHPSLVjgr6IYkpY9Aqqqid5.jpg",
			},
			{
				"id":         53,
				"name":       "Thriller",
				"posterPath": "https://image.tmdb.org/t/p/w500/3k8jv1kSAAc0rCfFGtWDDQL4dfK.jpg",
			},
			{
				"id":         10752,
				"name":       "War",
				"posterPath": "https://image.tmdb.org/t/p/w500/1v5ZteB49M0RUGYrf9R37Mz8yo2.jpg",
			},
			{
				"id":         37,
				"name":       "Western",
				"posterPath": "https://image.tmdb.org/t/p/w500/cvaBVpS0GzKqBd63pFT8f1E8OKv.jpg",
			},
		}
		return genres, nil
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

func (h *MovieHandler) FetchRecommendationsGrouped() AppHandler {
	return func(w http.ResponseWriter, r *http.Request) (interface{}, error) {
		var req struct {
			Titles []string `json:"titles"`
		}

		// Decode the JSON request body into the struct
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			return nil, err
		}
		return h.Service.GetRecommendationsGrouped(req.Titles)
	}
}

func (h *MovieHandler) SearchMovies() AppHandler {
	return func(w http.ResponseWriter, r *http.Request) (interface{}, error) {

		// Extract movie name from query parameters as integer
		movieName := r.URL.Query().Get("movieName")
		if movieName == "" {
			return nil, nil
		}
		// Use the extracted movie name to get recommendations
		return h.Service.SearchMovie(movieName)
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

func getIntFromRequest(r *http.Request, paramName string) int {
	param := r.URL.Query().Get(paramName)
	paramInt := 0
	if param != "" {
		if parsedParam, err := strconv.Atoi(param); err == nil {
			paramInt = parsedParam
		}
	}
	return paramInt
}
