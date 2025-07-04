package middleware

import (
	"context"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"net/http"
	"os"
	"time"
)

type Err error
type contextKey string

const UserIDKey contextKey = "userID"

func VerifyToken(r *http.Request) (int, error) {
	cookie, err := r.Cookie("token")
	if err != nil {
		log.Println("error to parse cookie", err)
		return 0, err
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
		log.Println("invalid token:", err)
		return 0, err
	}

	var userId float64
	claims, ok := token.Claims.(*MyCustomClaims)
	if !ok || !token.Valid {
		log.Println("invalid token claims:", err)
		return 0, err
	}

	if time.Now().Unix() > claims.ExpiresAt.Unix() {
		log.Println("token expired:", err)
		return 0, err
	}

	userId = claims.UserId
	log.Println("[VERIFY] Good token. UserID -", userId)
	return int(userId), nil
}

func WithAuth(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId, err := VerifyToken(r)
		if err != nil {
			http.Error(w, "failed to verify token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, userId)
		handler(w, r.WithContext(ctx))
	}
}
