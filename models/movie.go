package models

// Movie model
type Movie struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Year   int    `json:"year"`
	Poster string `json:"poster"`
}
