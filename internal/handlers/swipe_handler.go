package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

type Swipe struct {
	UserId int    `json:"user_id"`
	FilmId int    `json:"film_id"`
	Action string `json:"action"`
}

func SwipeHandler(w http.ResponseWriter, r *http.Request) {
	var swipe Swipe

	err := json.NewDecoder(r.Body).Decode(&swipe)
	if err != nil {
		log.Println("error to decode JSON")
		return
	}

}
