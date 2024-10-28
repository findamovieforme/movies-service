package models

// User model
type User struct {
	ID          int          `json:"id"`
	Name        string       `json:"name"`
	Email       string       `json:"email"`
	Preferences []Preference `json:"preferences"`
}

// Preference model
type Preference struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
