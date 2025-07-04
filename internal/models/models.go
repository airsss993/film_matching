package models

import "time"

type User struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	Password     string    `json:"password"`
	RegisterTime time.Time `json:"register_time"`
	//Role           string
	//FavoriteGenres []string
	//AvatarURL      string
	//LastLogin      time.Time
}

type Film struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	ReleaseYear int       `json:"release_year"`
	PosterURL   string    `json:"poster_url"`
	Rating      float32   `json:"rating"`
	Genres      []string  `json:"genres"`
	CreatedAt   time.Time `json:"created_at"`
	//RatingCount int
	//Country     string
	//DurationMin int
	//Director    string
	//ReleaseDate time.Time
}
