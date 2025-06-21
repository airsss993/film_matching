package handlers

import (
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"log"
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
	}

	if len(user.Name) <= 3 || len(user.Password) <= 5 {
		http.Error(w, "field is too short", http.StatusBadRequest)
	}

	// добавить обработку уже существующего пользователя по email

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	if err != nil {
		log.Println("failed to hash password", err)
	}

	user.Password = string(hash)

	service.AddUserToDB(user.Name, user.Email, user.Password)
}
