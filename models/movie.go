package models

import "github.com/ryanbradynd05/go-tmdb"

// Movie model
type Movie = tmdb.MovieShort

type Genre struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// GenreType represents a fixed set of GenreTypes
type GenreType string

const (
	Action      GenreType = "Action"
	Adventure   GenreType = "Adventure"
	Animation   GenreType = "Animation"
	Biography   GenreType = "Biography"
	Comedy      GenreType = "Comedy"
	Crime       GenreType = "Crime"
	Documentary GenreType = "Documentary"
	Drama       GenreType = "Drama"
	Family      GenreType = "Family"
	Fantasy     GenreType = "Fantasy"
	History     GenreType = "History"
	Horror      GenreType = "Horror"
	Music       GenreType = "Music"
	Musical     GenreType = "Musical"
	Mystery     GenreType = "Mystery"
	Romance     GenreType = "Romance"
	SciFi       GenreType = "Sci-Fi"
	Sport       GenreType = "Sport"
	Thriller    GenreType = "Thriller"
	War         GenreType = "War"
	Western     GenreType = "Western"
	Unknown     GenreType = "Unknown"
)

// Returned from tmdb
// [
//     "Action",
//     "Adventure",
//     "Animation",
//     "Comedy",
//     "Crime",
//     "Documentary",
//     "Drama",
//     "Family",
//     "Fantasy",
//     "History",
//     "Horror",
//     "Music",
//     "Mystery",
//     "Romance",
//     "Science Fiction",
//     "TV Movie",
//     "Thriller",
//     "War",
//     "Western"
// ]
