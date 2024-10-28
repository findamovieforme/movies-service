package services

import "github.com/movierecuh/movies-service/models"

type RecommenderServiceInterface interface {
	GetRecommendations(userID string) ([]models.Recommendation, error)
}

type RecommenderService struct{}

func (s *RecommenderService) GetRecommendations(userID string) ([]models.Recommendation, error) {
	// Simulate fetching user and recommendations
	user := models.User{
		ID:    1,
		Name:  "User 1",
		Email: "",
		Preferences: []models.Preference{
			{ID: 1, Name: "Preference 1"},
			{ID: 2, Name: "Preference 3"},
		},
	}

	movies := []models.Movie{
		{ID: 1, Title: "Movie 111", Year: 2020, Poster: ""},
		{ID: 2, Title: "Movie 2", Year: 2021, Poster: ""},
	}

	recommendations := []models.Recommendation{
		{ID: 1, User: user, Movie: movies[0]},
		{ID: 2, User: user, Movie: movies[1]},
	}

	return recommendations, nil
}
