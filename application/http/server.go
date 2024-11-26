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

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func setupRouter(env ServerEnv) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.Use(CORSMiddleware())

	router.Static("/thumbs", "./uploads")

	users := router.Group("/user")
	users.POST("/login", env.login)
	users.POST("/signup", env.signup)
	users.GET("/:id", env.getUserProfile)
	users.DELETE("/delete/:id", env.deleteUser)
	users.PATCH("/update/:id", env.updateUser)
	users.PATCH("/follow/:user_id", env.followUser)
	users.PATCH("/unfollow/:user_id", env.unfollowUser)

	// Pode ser removida mais tarde, apenas auxiliar
	users.GET("/all", env.getAllUsers)

	streams := router.Group("/livestreams")
	streams.POST("/create", env.createLiveStream)
	streams.DELETE("/delete/:id", env.deleteLiveStream)
	streams.PATCH("/update/:id", env.updateLiveStream)
	streams.GET("/feed", env.getFeed)
	streams.GET("/:user_id", env.getUserLiveStreams)
	streams.GET("/info/:id", env.getLiveStreamData)
	streams.GET("/on_publish", env.validateStream)

	streams.GET("/all", env.getAllStreams)

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
