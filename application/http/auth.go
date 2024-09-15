package http

import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
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

// Middleware que verifica se o token gerado por uma operação de
// login ou signup está no header `Authorization`. Se sim, proceder
// com a requisição.
func authMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		auth := ctx.Request.Header.Get("Authorization")
		if auth == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "missing Authorization header"})
			return
		}

		words := strings.Split(auth, " ")
		if len(words) != 2 {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "token not present in header"})
			return
		}

		token, err := jwt.Parse(words[1], func(t *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("ACCESS_TOKEN_SECRET")), nil
		})

		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "could not verify token " + err.Error()})
			return
		}

		if !token.Valid {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "invalid token"})
			return
		}

		ctx.Next()
	}
}
