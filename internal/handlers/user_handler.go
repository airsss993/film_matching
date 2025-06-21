package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log"
	"main/internal/db"
	"main/internal/service"
	"net/http"
)

func RegistrationHandler(w http.ResponseWriter, r *http.Request) {
	type User struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var user User

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
	// добавить обработку уже существующего пользователя по email

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	if err != nil {
		log.Println("failed to hash password", err)
		return
	}

	user.Password = string(hash)

	service.AddUserToDB(user.Name, user.Email, user.Password)
}
