// Package classification LiveStreamAPI
//
// Documentação da API de liveStreams
//
//	Schemes: http
//	Version: 1.0.0
//
//	Consumes:
//	- application/json
//
// swagger:meta
package http

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gtvb/livestream/models"
)

type ServerEnv struct {
	liveStreamsRepository models.LiveStreamRepositoryInterface
	userRepository        models.UserRepositoryInterface
}

// Inicia um servidor HTTP e define as rotas padrão da aplicação
func RunServer(lr models.LiveStreamRepositoryInterface, ur models.UserRepositoryInterface) {
	env := ServerEnv{
		liveStreamsRepository: lr,
		userRepository:        ur,
	}

	router := gin.Default()

	users := router.Group("/user")
	users.POST("/login", env.login)
	users.POST("/signup", env.signup)
	users.GET("/all", env.getAllUsers)
	users.DELETE("/delete", authMiddleware(), env.deleteUser)
	users.PATCH("/update", authMiddleware(), env.updateUser)
	users.GET("/:id", authMiddleware(), env.getUserProfile)

	streams := router.Group("/livestream")
	streams.POST("/create", authMiddleware(), env.createLiveStream)
	streams.DELETE("/delete", authMiddleware(), env.deleteLiveStream)
	streams.PATCH("/update", authMiddleware(), env.updateLiveStream)
	streams.GET("/:user_id", authMiddleware(), env.getUserLiveStreams)

	streams.POST("/publish_auth", func(ctx *gin.Context) {
		var body interface{}

		if err := ctx.ShouldBindBodyWithJSON(&body); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "could not parse body"})
			return
		}

		fmt.Println(body)

		ctx.JSON(http.StatusOK, gin.H{"message": "allowed to proceed"})
	})

	router.Run(":" + os.Getenv("SERVER_PORT"))
}
