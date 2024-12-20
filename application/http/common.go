package http

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/gtvb/livestream/infra/db"
	"github.com/gtvb/livestream/infra/repository"
	"github.com/gtvb/livestream/utils"
)

func setupDatabase() *utils.TestContainer {
	container, err := utils.NewTestContainer("ls-db-test")
	if err != nil {
		log.Panicf("Error: could not start container, reason -> %s\n", err)
	}

	err = container.SetupDatabaseWrapper()
	if err != nil {
		log.Panicf("Error: could not start database wrapper, reason -> %s\n", err)
	}

	return container
}

func setupEnv(database *db.Database) ServerEnv {
	env := ServerEnv{}

	userRepo := repository.NewUserRepository(database, "users_test")
	liveStreamRepo := repository.NewLiveStreamRepository(database, "livestreams_test")

	env.userRepository = userRepo
	env.liveStreamsRepository = liveStreamRepo

	return env
}

func makeRequest(router *gin.Engine, method, url string, body interface{}) *httptest.ResponseRecorder {
	requestBody, _ := json.Marshal(body)
	req, _ := http.NewRequest(method, url, bytes.NewBuffer(requestBody))

	writer := httptest.NewRecorder()
	router.ServeHTTP(writer, req)

	return writer
}
