package models

import "github.com/ryanbradynd05/go-tmdb"

// Movie model
type Movie = tmdb.MovieShort

// Genre represents a fixed set of genres
type Genre string

const (
	Action      Genre = "Action"
	Adventure   Genre = "Adventure"
	Animation   Genre = "Animation"
	Biography   Genre = "Biography"
	Comedy      Genre = "Comedy"
	Crime       Genre = "Crime"
	Documentary Genre = "Documentary"
	Drama       Genre = "Drama"
	Family      Genre = "Family"
	Fantasy     Genre = "Fantasy"
	History     Genre = "History"
	Horror      Genre = "Horror"
	Music       Genre = "Music"
	Musical     Genre = "Musical"
	Mystery     Genre = "Mystery"
	Romance     Genre = "Romance"
	SciFi       Genre = "Sci-Fi"
	Sport       Genre = "Sport"
	Thriller    Genre = "Thriller"
	War         Genre = "War"
	Western     Genre = "Western"
	Unknown     Genre = "Unknown"
)
