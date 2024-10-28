package models

// Recommendation model
type Recommendation struct {
	ID    int   `json:"id"`
	User  User  `json:"user"`
	Movie Movie `json:"movie"`
}
