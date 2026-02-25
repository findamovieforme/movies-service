package routers

import (
	"net/http"
	"sync"
	"time"

	"github.com/findamovieforme/movies-service/handlers"
	"github.com/findamovieforme/movies-service/services"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type ipRateLimiter struct {
	mu       sync.Mutex
	requests map[string][]time.Time
	limit    int
	window   time.Duration
}

func newIPRateLimiter(limit int, window time.Duration) *ipRateLimiter {
	return &ipRateLimiter{
		requests: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
	}
}

func (l *ipRateLimiter) Allow(ip string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-l.window)

	times := l.requests[ip]
	filtered := times[:0]
	for _, t := range times {
		if t.After(cutoff) {
			filtered = append(filtered, t)
		}
	}

	if len(filtered) >= l.limit {
		l.requests[ip] = filtered
		return false
	}

	filtered = append(filtered, now)
	l.requests[ip] = filtered
	return true
}

func rateLimitMiddleware(l *ipRateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		if !l.Allow(ip) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "rate limit exceeded",
			})
			return
		}
		c.Next()
	}
}

func InitRouter() *gin.Engine {
	router := gin.Default()
	router.Use(cors.Default())

	// Initialize services and handlers
	movieService := services.GetMovieService()

	defaultHandler := handlers.NewDefaultHandler()
	movieHandler := handlers.NewMovieHandler(movieService)

	// Per-IP rate limiter specifically for expensive GPT endpoint
	gptLimiter := newIPRateLimiter(20, time.Minute)

	// Add a base /movies to all routes
	moviesGroup := router.Group("/movies")

	// Define routes and their corresponding handlers
	moviesGroup.GET("/", defaultHandler.Default)
	moviesGroup.GET("/ping", defaultHandler.Ping)
	moviesGroup.GET("/trending", gin.WrapF(movieHandler.FetchTrendingMovies().ServeHTTP))
	moviesGroup.GET("/recent", gin.WrapF(movieHandler.FetchRecentlyReleasedMovies().ServeHTTP))
	moviesGroup.GET("/genres", gin.WrapF(movieHandler.FetchGenres().ServeHTTP))
	moviesGroup.GET("/details", gin.WrapF(movieHandler.FetchMovieDetails().ServeHTTP))
	moviesGroup.GET("/search", gin.WrapF(movieHandler.SearchMovies().ServeHTTP))

	moviesGroup.POST("/recommendations", gin.WrapF(movieHandler.FetchRecommendations().ServeHTTP))
	moviesGroup.POST("/recommendationsGrouped", gin.WrapF(movieHandler.FetchRecommendationsGrouped().ServeHTTP))
	moviesGroup.POST("/gpt", rateLimitMiddleware(gptLimiter), gin.WrapF(movieHandler.FetchGptResponse().ServeHTTP))

	return router
}
