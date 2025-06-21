package handlers

import (
	"golang.org/x/crypto/bcrypt"
	"log"
	"main/internal/service"
	"net/http"
)

func UserRoute(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	email := r.URL.Query().Get("email")
	password := r.URL.Query().Get("password")

	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		log.Println("failed to hash password", err)
	}
	password = string(hash)

	service.AddUser(name, email, password)
	w.Write([]byte("User successfully created!"))
}
