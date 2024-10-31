package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/movierecuh/movies-service/handlers"
	"github.com/movierecuh/movies-service/services"
)

func InitRouter() *gin.Engine {
	router := gin.Default()

	// Initialize services and handlers
	recommenderService := &services.RecommenderService{}
	recommenderHandler := handlers.NewRecommenderHandler(recommenderService)

	movieService := services.GetMovieService()
	movieHandler := handlers.NewMovieHandler(movieService)

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// Define routes and their corresponding handlers
	router.GET("/recommendations", gin.WrapF(recommenderHandler.FetchRecommendations))
	router.GET("/movies", gin.WrapF(movieHandler.Fetchmovies().ServeHTTP))

	return router
}
