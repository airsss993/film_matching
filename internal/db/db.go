package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func ConnDB() *sql.DB {
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "localhost"
	}

	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "5432"
	}

	sslMode := os.Getenv("DB_SSL_MODE")
	if sslMode == "" {
		sslMode = "disable"
	}

	connStr := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=%s",
		host, port, os.Getenv("DB_USER"), os.Getenv("DB_NAME"), sslMode)

	log.Printf("Connecting to DB with: %s", connStr)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	// Проверяем соединение
	if err = db.Ping(); err != nil {
		log.Printf("Failed to ping database: %v", err)
		log.Fatal(err)
	}

	log.Println("Successfully connected to database")
	return db
}
