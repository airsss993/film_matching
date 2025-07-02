package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type TMDBMovie struct {
	ID          int     `json:"id"`
	Title       string  `json:"title"`
	Overview    string  `json:"overview"`
	ReleaseDate string  `json:"release_date"`
	PosterPath  string  `json:"poster_path"`
	GenreIDs    []int   `json:"genre_ids"`
	VoteAverage float64 `json:"vote_average"`
}

type TMDBResponse struct {
	Results []TMDBMovie `json:"results"`
}

var genreMap = map[int]string{
	28: "Action", 12: "Adventure", 16: "Animation", 35: "Comedy",
	80: "Crime", 99: "Documentary", 18: "Drama", 10751: "Family",
	14: "Fantasy", 36: "History", 27: "Horror", 10402: "Music",
	9648: "Mystery", 10749: "Romance", 878: "Sci-Fi", 10770: "TV Movie",
	53: "Thriller", 10752: "War", 37: "Western",
}

func main() {
	_ = godotenv.Load()

	apiKey := os.Getenv("TMDB_API_KEY")
	dbURL := os.Getenv("DB_URL")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("DB connect error: %v", err)
	}
	defer db.Close()

	const pages = 3

	for page := 1; page <= pages; page++ {
		url := fmt.Sprintf("https://api.themoviedb.org/3/movie/popular?api_key=%s&page=%d", apiKey, page)
		resp, err := http.Get(url)
		if err != nil {
			log.Fatalf("TMDB error: %v", err)
		}
		defer resp.Body.Close()

		var data TMDBResponse
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			log.Fatalf("JSON error: %v", err)
		}

		for _, movie := range data.Results {
			if movie.Title == "" {
				continue
			}

			year := 0
			if len(movie.ReleaseDate) >= 4 {
				fmt.Sscanf(movie.ReleaseDate[:4], "%d", &year)
			}

			genres := []string{}
			for _, id := range movie.GenreIDs {
				if name, ok := genreMap[id]; ok {
					genres = append(genres, name)
				}
			}

			posterURL := ""
			if movie.PosterPath != "" {
				posterURL = "https://image.tmdb.org/t/p/w500" + movie.PosterPath
			}

			var exists bool
			err := db.QueryRow("SELECT EXISTS (SELECT 1 FROM films WHERE title = $1 AND release_year = $2)", movie.Title, year).Scan(&exists)
			if err != nil {
				log.Printf("DB check error: %v", err)
				continue
			}
			if exists {
				continue
			}

			_, err = db.ExecContext(context.Background(), `
				INSERT INTO films (title, description, release_year, poster_url, rating, genres)
				VALUES ($1, $2, $3, $4, $5, $6)
			`, movie.Title, movie.Overview, year, posterURL, movie.VoteAverage, strings.Join(genres, ", "))
			if err != nil {
				log.Printf("insert error: %v", err)
			} else {
				fmt.Printf("Add film: %s (%d)\n", movie.Title, year)
			}
		}
	}
}
