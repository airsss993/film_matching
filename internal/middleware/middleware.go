package middleware

import (
	"github.com/golang-jwt/jwt/v5"
	"log"
	"net/http"
	"os"
	"time"
)

func VerifyToken(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("token")
	if err != nil {
		http.Error(w, "error to parse cookie", http.StatusUnauthorized)
		return
	}

	tokenString := cookie.Value

	type MyCustomClaims struct {
		UserId float64 `json:"user_id"`
		jwt.RegisteredClaims
	}

	token, err := jwt.ParseWithClaims(tokenString, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET")), nil
	}, jwt.WithLeeway(5*time.Second))
	if err != nil {
		http.Error(w, "Token expired", http.StatusUnauthorized)
		return
	}

	if claims, ok := token.Claims.(*MyCustomClaims); ok {
		if claims.ExpiresAt != nil && time.Now().After(claims.ExpiresAt.Time) {
			http.Error(w, "Token expired", http.StatusUnauthorized)
			return
		} else {
			log.Println("Good token.")
		}
	}
}
