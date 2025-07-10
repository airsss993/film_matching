package handlers

import (
	"encoding/json"
	"log"
	"main/internal/db"
	"main/internal/middleware"
	"net/http"
)

type SwipeInput struct {
	FilmID int  `json:"film_id"`
	Liked  bool `json:"liked"`
}

type SwipeResponse struct {
	IsMatch bool `json:"is_match"`
}

func SwipeHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey)
	if userID == nil {
		http.Error(w, "user ID not found in context", http.StatusInternalServerError)
		return
	}

	var swipe SwipeInput
	err := json.NewDecoder(r.Body).Decode(&swipe)
	if err != nil {
		log.Println("error to decode JSON:", err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	DB := db.ConnDB()

	action := "skip"
	if swipe.Liked {
		action = "like"
	}

	_, err = DB.Exec(`INSERT INTO swipes (user_id, film_id, action) VALUES ($1, $2, $3)`, userID, swipe.FilmID, action)
	if err != nil {
		log.Println("failed to insert swipe:", err)
		http.Error(w, "failed to save swipe", http.StatusInternalServerError)
		return
	}

	log.Printf("User %v swiped %s on film %v", userID, action, swipe.FilmID)

	isMatch := false

	response := SwipeResponse{IsMatch: isMatch}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		log.Println("failed to encode response:", err)
		return
	}
}
