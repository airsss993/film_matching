package main

import (
	"log"
	"main/internal/db"
	"main/internal/handlers"
	"main/internal/middleware"
	"net/http"

	"github.com/joho/godotenv"
)

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
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

	db.ConnDB()

	http.Handle("/register", enableCORS(http.HandlerFunc(handlers.RegistrationHandler)))
	http.Handle("/login", enableCORS(http.HandlerFunc(handlers.LoginHandler)))
	http.Handle("/api/films/next", enableCORS(middleware.WithAuth(handlers.GetNextFilm)))
	http.Handle("/api/films/swipe", enableCORS(middleware.WithAuth(handlers.SwipeHandler)))
	http.Handle("/api/matches", enableCORS(middleware.WithAuth(handlers.GetMatches)))
	http.Handle("/api/users/profile", enableCORS(middleware.WithAuth(handlers.GetUserProfile)))
	http.Handle("/api/auth/logout", enableCORS(http.HandlerFunc(handlers.LogoutHandler)))

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("../../web/static"))))
	http.HandleFunc("/auth.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../../web/auth.html")
	})
	http.HandleFunc("/home.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../../web/home.html")
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.Redirect(w, r, "/auth.html", http.StatusFound)
			return
		}
		http.NotFound(w, r)
	})

	log.Println("Server starts at http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("error starting server: ", err)
	}
}
