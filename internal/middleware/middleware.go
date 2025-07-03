package middleware

import (
	"github.com/golang-jwt/jwt/v5"
	"log"
	"net/http"
	"os"
	"time"
)

type Err error

func VerifyToken(w http.ResponseWriter, r *http.Request) {
	func() (float64, Err) {
		cookie, err := r.Cookie("token")
		if err != nil {
			http.Error(w, "error to parse cookie", http.StatusUnauthorized)
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
			http.Error(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
			return 0, err
		}

		var userId float64
		claims, ok := token.Claims.(*MyCustomClaims)
		if !ok || !token.Valid {
			http.Error(w, "invalid token claims", http.StatusUnauthorized)
			return 0, err
		}

		if time.Now().Unix() > claims.ExpiresAt.Unix() {
			http.Error(w, "token expired", http.StatusUnauthorized)
			return 0, err
		}

		userId = claims.UserId
		log.Println("[VERIFY] Good token. UserID -", userId)
		return userId, nil
	}()
}

//func WithAuth(w http.ResponseWriter, r *http.Request) http.HandlerFunc {
//	token, err := VerifyToken(r.Cookie("token"))
//	if err != nil {
//		log.Println("[AUTH] Token not found.")
//	}
//
//}
