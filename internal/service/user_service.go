package service

import (
	"log"
	"main/internal/db"
)

func AddUser(name, email, password string) {
	DB := db.ConnDB()

	_, err := DB.Exec(`INSERT INTO users(username, email, passwordhash) VALUES ($1, $2, $3)`, name, email, password)
	if err != nil {
		log.Println("error to insert user", err)
	} else {
		log.Println("User successfully added to DB!")
	}
}
