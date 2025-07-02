package main

import (
	"github.com/joho/godotenv"
	"log"
	"main/internal/db"
	"main/internal/handlers"
	"net/http"
)

func main() {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Print("no .env file found")
	}

	http.HandleFunc("/auth/register", handlers.RegistrationHandler)
	http.HandleFunc("/auth/login", handlers.LoginHandler)

	log.Println("Server starts at http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("error starting server")
	}

	db.ConnDB()
}
