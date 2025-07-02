package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"log"
	"main/internal/db"
	"main/internal/service"
	"net/http"
	"os"
	"time"
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
	fmt.Println("USER ID FROM DB", userId)
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
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}

	http.SetCookie(w, cookie)
	log.Println("Token generated.")
}
