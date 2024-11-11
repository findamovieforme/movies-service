package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/movierecuh/movies-service/handlers"
	"github.com/movierecuh/movies-service/services"
)

func InitRouter() *gin.Engine {
	router := gin.Default()

	// Initialize services and handlers
	movieService := services.GetMovieService()

	defaultHandler := handlers.NewDefaultHandler()
	movieHandler := handlers.NewMovieHandler(movieService)
	// Add a base /movies to all routes
	moviesGroup := router.Group("/movies")

	// Define routes and their corresponding handlers
	moviesGroup.GET("/", defaultHandler.Default)
	moviesGroup.GET("/ping", defaultHandler.Ping)
	moviesGroup.GET("/trending", gin.WrapF(movieHandler.FetchTrendingMovies().ServeHTTP))
	moviesGroup.GET("/recent", gin.WrapF(movieHandler.FetchRecentlyReleasedMovies().ServeHTTP))
	moviesGroup.GET("/genres", gin.WrapF(movieHandler.FetchGenres().ServeHTTP))
	moviesGroup.GET("/details", gin.WrapF(movieHandler.FetchMovieDetails().ServeHTTP))

	moviesGroup.POST("/recommendations", gin.WrapF(movieHandler.FetchRecommendations().ServeHTTP))

	return router
}
