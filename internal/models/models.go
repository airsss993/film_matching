package models

import "time"

type User struct {
	ID           int
	Name         string
	Email        string
	PasswordHash string
	RegisteredAt time.Time
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
