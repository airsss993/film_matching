package handlers

import (
	"main/internal/service"
	"net/http"
)

func UserRoute(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	email := r.URL.Query().Get("email")
	password := r.URL.Query().Get("password")

	service.AddUser(name, email, password)
	w.Write([]byte("User successfully created!"))
}
