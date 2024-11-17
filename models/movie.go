package models

import "github.com/ryanbradynd05/go-tmdb"

// Movie model
type Movie = tmdb.MovieShort

type ExtendedMovie struct {
	Adult         bool    `json:"adult"`
	BackdropPath  string  `json:"backdrop_path"`
	ID            int     `json:"id"`
	OriginalTitle string  `json:"original_title"`
	GenreIDs      []int32 `json:"genre_ids"`
	Popularity    float32 `json:"popularity"`
	PosterPath    string  `json:"poster_path"`
	ReleaseDate   string  `json:"release_date"`
	Title         string  `json:"title"`
	Overview      string  `json:"overview"`
	Video         bool    `json:"video"`
	VoteAverage   float32 `json:"vote_average"`
	VoteCount     uint32  `json:"vote_count"`
	TrailerKey    string  `json:"trailer_key"`
}

type Genre struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	PosterPath string `json:"posterPath"`
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
