package http

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	UserId int `json:"user_id"`
	jwt.RegisteredClaims
}

// Gera um token JWT utilizando o id de usuário como claim.
// Para saber mais sobre claims visite jwt.io
func generateToken(userId interface{}) (string, error) {
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
