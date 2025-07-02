package main

import (
	"github.com/joho/godotenv"
	"log"
	"main/internal/db"
	"main/internal/handlers"
	"main/internal/middleware"
	"net/http"
)

func main() {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Print("no .env file found")
	}

	http.HandleFunc("/register", handlers.RegistrationHandler)
	http.HandleFunc("/login", handlers.LoginHandler)
	http.HandleFunc("/auth", middleware.VerifyToken)

	log.Println("Server starts at http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("error starting server")
	}

	db.ConnDB()
}
