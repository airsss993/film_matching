package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"main/internal/db"
	"main/internal/service"
	"net/http"
	"os"
	"time"

	"main/internal/middleware"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterInput struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func RegistrationHandler(w http.ResponseWriter, r *http.Request) {
	var user RegisterInput

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Println("error to decode JSON")
		return
	}

	if len(user.Name) <= 3 || len(user.Password) <= 5 {
		http.Error(w, "field is too short", http.StatusBadRequest)
		return
	}

	var existingEmail string

	DB := db.ConnDB()
	err = DB.QueryRow(`SELECT email FROM users WHERE email = $1`, user.Email).Scan(&existingEmail)
	fmt.Println(existingEmail)
	if err == nil {
		http.Error(w, "user with this email already exists", http.StatusBadRequest)
		return
	}
	if err != sql.ErrNoRows {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	if err != nil {
		log.Println("failed to hash password", err)
		return
	}

	user.Password = string(hash)

	service.AddUserToDB(user.Name, user.Email, user.Password)

	log.Println("Successfully registered new account!")
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	var user LoginInput

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Println("error to decode JSON")
		return
	}

	DB := db.ConnDB()

	var userId int
	var hashedPassword string
	err = DB.QueryRow(`SELECT id, password FROM users WHERE email = $1`, user.Email).Scan(&userId, &hashedPassword)
	if err != nil {
		http.Error(w, "failed to get id or password", http.StatusBadRequest)
		log.Println(err)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(user.Password))
	if err != nil {
		http.Error(w, "failed to compare password", http.StatusBadRequest)
		return
	}
	log.Println("User data verification: matched.")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": float64(userId),
		"exp":     time.Now().Add(900 * time.Second).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		log.Println("error to sign token", err)
		return
	}

	cookie := &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: false,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, cookie)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Login successful",
	})

	log.Println("Token generated.")
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	//w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	//w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	//w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	//w.Header().Set("Access-Control-Allow-Credentials", "true")
	//
	//if r.Method == "OPTIONS" {
	//	w.WriteHeader(http.StatusOK)
	//	return
	//}

	cookie := &http.Cookie{
		Name:     "token",
		Value:    "",
		Path:     "/",
		Expires:  time.Now().Add(-time.Hour),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, cookie)

	w.WriteHeader(http.StatusOK)
}

type UserProfile struct {
	Name         string   `json:"name"`
	TotalSwipes  int      `json:"total_swipes"`
	TotalMatches int      `json:"total_matches"`
	Preferences  []string `json:"preferences"`
}

type Match struct {
	Film struct {
		ID        int     `json:"id"`
		Title     string  `json:"title"`
		Year      int     `json:"year"`
		Rating    float32 `json:"rating"`
		PosterURL string  `json:"poster_url"`
	} `json:"film"`
}

func GetUserProfile(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey)
	if userID == nil {
		http.Error(w, "user ID not found in context", http.StatusInternalServerError)
		return
	}

	DB := db.ConnDB()

	var profile UserProfile

	err := DB.QueryRow(`SELECT username FROM users WHERE id = $1`, userID).Scan(&profile.Name)
	if err != nil {
		log.Println("failed to get user name:", err)
		http.Error(w, "failed to get user profile", http.StatusInternalServerError)
		return
	}

	err = DB.QueryRow(`SELECT COUNT(*) FROM swipes WHERE user_id = $1`, userID).Scan(&profile.TotalSwipes)
	if err != nil {
		log.Println("failed to get swipes count:", err)
		profile.TotalSwipes = 0
	}

	profile.TotalMatches = 0
	profile.Preferences = []string{"Action", "Drama", "Comedy"}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(profile); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		log.Println("failed to encode response:", err)
		return
	}
}

func GetMatches(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey)
	if userID == nil {
		http.Error(w, "user ID not found in context", http.StatusInternalServerError)
		return
	}

	matches := []Match{}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(matches); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		log.Println("failed to encode response:", err)
		return
	}
}
