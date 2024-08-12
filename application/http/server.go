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
	"net/http"
	"os"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
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

	router := chi.NewRouter()
	router.Use(middleware.Logger)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello Go!"))
	})

	router.Route("/users", func(router chi.Router) {
		router.Post("/login", env.login)
		router.Post("/signup", env.signup)
	})

	http.ListenAndServe(":"+os.Getenv("SERVER_PORT"), router)
}
