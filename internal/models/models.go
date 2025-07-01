package models

import "time"

type User struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	Password     string    `json:"password"`
	RegisterTime time.Time `json:"register_time"`
	//Role           string
	//FavoriteGenres []string
	//AvatarURL      string
	//LastLogin      time.Time
}

type Film struct {
	ID          int
	Title       string
	Description string
	Year        int
	//Genres      []string
	//Rating      float32
	//RatingCount int
	//PosterURL   string
	//Country     string
	//DurationMin int
	//Director    string
	//ReleaseDate time.Time
}
