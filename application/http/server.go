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
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gtvb/livestream/models"
)

type ServerEnv struct {
	liveStreamsRepository models.LiveStreamRepositoryInterface
	userRepository        models.UserRepositoryInterface
}

func setupRouter(env ServerEnv) *gin.Engine {
	router := gin.Default()

	users := router.Group("/user")
	users.POST("/login", env.login)
	users.POST("/signup", env.signup)
	users.GET("/:id", authMiddleware(), env.getUserProfile)
	users.DELETE("/delete/:id", authMiddleware(), env.deleteUser)
	users.PATCH("/update/:id", authMiddleware(), env.updateUser)

	users.GET("/all", env.getAllUsers)

	streams := router.Group("/livestreams")
	streams.POST("/create", authMiddleware(), env.createLiveStream)
	streams.DELETE("/delete/:id", authMiddleware(), env.deleteLiveStream)
	streams.PATCH("/update/:id", authMiddleware(), env.updateLiveStream)
	streams.GET("/:user_id", authMiddleware(), env.getUserLiveStreams)
	streams.GET("/info/:id", authMiddleware(), env.getLiveStreamData)
	streams.GET("/on_publish", env.validateStream)

	return router
}

// Inicia um servidor HTTP e define as rotas padrão da aplicação
func RunServer(lr models.LiveStreamRepositoryInterface, ur models.UserRepositoryInterface) {
	env := ServerEnv{
		liveStreamsRepository: lr,
		userRepository:        ur,
	}

	router := setupRouter(env)
	router.Run(":" + os.Getenv("SERVER_PORT"))
}
