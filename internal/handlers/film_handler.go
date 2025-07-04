package handlers

import (
	"encoding/json"
	"log"
	"main/internal/db"
	"main/internal/middleware"
	"main/internal/models"
	"net/http"
)

func GetFilm(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey)
	if userID == nil {
		http.Error(w, "user ID not found in context", http.StatusInternalServerError)
		return
	}
	log.Println("UserID -", userID)

	var film models.Film

	DB := db.ConnDB()
	row := DB.QueryRow(`
		SELECT id, title, description, release_year, poster_url, rating, genres, created_at 
		FROM films 
		WHERE id = 1`)
	if row == nil {
		log.Println("failed to get data")
	}

	var genres string
	err := row.Scan(
		&film.ID,
		&film.Title,
		&film.Description,
		&film.ReleaseYear,
		&film.PosterURL,
		&film.Rating,
		&genres,
		&film.CreatedAt,
	)
	if err != nil {
		log.Println("failed to scan data")
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(film); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		log.Println("failed to encode response:", err)
		return
	}
}
