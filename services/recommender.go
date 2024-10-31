package services

import (
	"fmt"

	"github.com/movierecuh/movies-service/models"
)

type RecommenderServiceInterface interface {
	GetRecommendations(userID string) ([]models.Movie, error)
}

type RecommenderService struct{}

func (s *RecommenderService) GetRecommendations(userID string) ([]models.Movie, error) {
	// UserID found
	if userID == "1" {
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

		fmt.Println("User found: ", user)
		movies := []models.Movie{
			{ID: 1, Title: "Movie 111", Year: 2020, Poster: ""},
			{ID: 2, Title: "Movie 2", Year: 2021, Poster: ""},
		}

		return movies, nil
	}

	// Return empty
	return []models.Movie{}, nil

}
