package handlers

import (
	"encoding/json"
	"log"
	"main/internal/db"
	"main/internal/middleware"
	"main/internal/models"
	"net/http"
)

func GetNextFilm(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey)
	if userID == nil {
		http.Error(w, "user ID not found in context", http.StatusInternalServerError)
		return
	}

	var film models.Film

	DB := db.ConnDB()
	row := DB.QueryRow(`SELECT * FROM films WHERE id NOT IN (SELECT film_id FROM swipes WHERE user_id = $1) LIMIT 1;`, userID)
	if row == nil {
		log.Println("failed to get data")
		http.Error(w, "no films found", http.StatusNotFound)
		return
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
		log.Println("failed to scan data:", err)
		http.Error(w, "failed to get film data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(film); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		log.Println("failed to encode response:", err)
		return
	}
}
