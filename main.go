package main

import (
	"context"
	"log"

	"github.com/gtvb/livestream/application/http"
	"github.com/gtvb/livestream/infra/db"
	"github.com/gtvb/livestream/infra/repository"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	db, err := db.NewDb()
	if err != nil {
		log.Printf("Failed to start database: %s\n", err.Error())
		return
	}
	defer func() {
		log.Printf("disconnecting from db\n")
		db.Client().Disconnect(context.TODO())
	}()

	userRepository := repository.NewUserRepository(db)
	liveStreamsRepository := repository.NewLiveStreamRepository(db)

	http.RunServer(liveStreamsRepository, userRepository)
}
