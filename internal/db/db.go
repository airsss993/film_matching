package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"
)

func ConnDB() *sql.DB {
	connStr := fmt.Sprintf("user=%s dbname=%s", os.Getenv("DB_USER"), os.Getenv("DB_NAME"))
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	//defer db.Close()

	return db
}
