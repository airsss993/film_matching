package main

import (
	"github.com/joho/godotenv"
	"log"
	"main/internal/db"
	"main/internal/handlers"
	"main/internal/middleware"
	"net/http"
)

// Добавьте эту функцию
func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Max-Age", "3600")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Print("no .env file found")
	}

	http.Handle("/register", enableCORS(http.HandlerFunc(handlers.RegistrationHandler)))
	http.Handle("/login", enableCORS(http.HandlerFunc(handlers.LoginHandler)))
	http.Handle("/film", enableCORS(middleware.WithAuth(handlers.GetFilm)))

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./web/static/"))))
	http.HandleFunc("/auth", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./web/auth.html")
	})

	log.Println("Server starts at http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("error starting server")
	}

	db.ConnDB()
}
