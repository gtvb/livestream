package http

import (
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Gera um token JWT utilizando o id de usuário como claim.
// Para saber mais sobre claims visite jwt.io
func generateToken(userId int) (string, error) {
	ttl := time.Hour * 5
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"user_id": userId,
			"exp":     time.Now().UTC().Add(ttl).Unix(),
		})

	signedToken, err := jwtToken.SignedString([]byte(os.Getenv("ACCESS_TOKEN_SECRET")))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

// Middleware que verifica se o token gerado por uma operação de
// login ou signup está no header `Authorization`. Se sim, proceder
// com a requisição.
func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		authorization := req.Header.Get("Authorization")
		if authorization == "" {
			http.Error(w, "no authorization header", http.StatusUnauthorized)
			return
		}

		// Auth header format: `Bearer <token>` <- token starts at index 7
		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(authorization[7:], claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("ACCESS_TOKEN_SECRET")), nil
		})
		if err != nil {
			http.Error(w, "failed to apply auth handler", http.StatusInternalServerError)
			return
		}

		if !token.Valid {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, req)
	}
}
