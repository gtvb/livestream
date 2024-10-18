package http

import (
	"net/http"
	"testing"

	"github.com/gtvb/livestream/models"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func createTestUser(env ServerEnv) *models.User {
	id, _ := env.userRepository.CreateUser("test_username", "test@email.com", hashPassword("test_pass"))
	userID := id.(primitive.ObjectID)

	return &models.User{
		ID:       userID,
		Username: "test_username",
		Email:    "test@email.com",
		Password: "test_pass",
	}
}

func TestCreateLiveStream(t *testing.T) {
	container := setupDatabase()
	defer container.Terminate()

	env := setupEnv(container.Database)
	user := createTestUser(env)

	authParams := &AuthParams{Email: user.Email, Password: user.Password}
	createLiveStreamBody := CreateLiveStreamBody{
		UserId: user.ID.Hex(),
		Name:   "Test Live Stream",
	}

	t.Run("Correct body", func(t *testing.T) {
		router := setupRouter(env)
		writer := makeRequest(router, "POST", "/livestreams/create", createLiveStreamBody, authParams)

		assert.Equal(t, http.StatusCreated, writer.Code)
		assert.Contains(t, writer.Body.String(), "livestream created")
	})

	t.Run("Invalid ID", func(t *testing.T) {
		createLiveStreamBody.UserId = "invalidID"
		router := setupRouter(env)
		writer := makeRequest(router, "POST", "/livestreams/create", createLiveStreamBody, authParams)

		assert.Equal(t, http.StatusBadRequest, writer.Code)
		assert.Contains(t, writer.Body.String(), "invalid ObjectID")
	})

	t.Run("No body on request", func(t *testing.T) {
		router := setupRouter(env)
		writer := makeRequest(router, "POST", "/livestreams/create", nil, authParams)

		assert.Equal(t, http.StatusBadRequest, writer.Code)
		assert.Contains(t, writer.Body.String(), "could not get body from request")
	})
}

func TestDeleteLiveStream(t *testing.T) {
	container := setupDatabase()
	defer container.Terminate()

	env := setupEnv(container.Database)
	user := createTestUser(env)
	router := setupRouter(env)
	authParams := &AuthParams{Email: user.Email, Password: user.Password}

	t.Run("Successfully create stream", func(t *testing.T) {
		streamID, err := env.liveStreamsRepository.CreateLiveStream("Test Stream", "streamkey-test", user.ID)
		if err != nil {
			t.Fatalf("Failed to create test live stream: %v", err)
		}
		id := streamID.(primitive.ObjectID)

		writer := makeRequest(router, "DELETE", "/livestreams/delete/"+id.Hex(), nil, authParams)
		assert.Equal(t, http.StatusOK, writer.Code)
		assert.Contains(t, writer.Body.String(), "success")
	})

	t.Run("Invalid ID", func(t *testing.T) {
		writer := makeRequest(router, "DELETE", "/livestreams/delete/invalidID", nil, authParams)
		assert.Equal(t, http.StatusBadRequest, writer.Code)
		assert.Contains(t, writer.Body.String(), "invalid ObjectID")
	})

	t.Run("Stream not found", func(t *testing.T) {
		writer := makeRequest(router, "DELETE", "/livestreams/delete/"+primitive.NewObjectID().Hex(), nil, authParams)
		assert.Equal(t, http.StatusNotFound, writer.Code)
		assert.Contains(t, writer.Body.String(), "failed to delete live stream")
	})
}

func TestGetLiveStreamInfo(t *testing.T) {
	container := setupDatabase()
	defer container.Terminate()

	env := setupEnv(container.Database)
	router := setupRouter(env)
	user := createTestUser(env)

	authParams := &AuthParams{Email: user.Email, Password: user.Password}

	streamID, err := env.liveStreamsRepository.CreateLiveStream("Test Stream", "streamkey-test", user.ID)
	if err != nil {
		t.Fatalf("Failed to create test live stream: %v", err)
	}
	id := streamID.(primitive.ObjectID)

	t.Run("Succesfully search for stream", func(t *testing.T) {
		writer := makeRequest(router, "GET", "/livestreams/info/"+id.Hex(), nil, authParams)
		assert.Equal(t, http.StatusOK, writer.Code)
		assert.Contains(t, writer.Body.String(), "Test Stream")
	})

	t.Run("Invalid ID", func(t *testing.T) {
		writer := makeRequest(router, "GET", "/livestreams/info/invalidID", nil, authParams)
		assert.Equal(t, http.StatusBadRequest, writer.Code)
		assert.Contains(t, writer.Body.String(), "invalid ObjectID")
	})

	t.Run("Unexistant ID", func(t *testing.T) {
		writer := makeRequest(router, "GET", "/livestreams/info/"+primitive.NewObjectID().Hex(), nil, authParams)
		assert.Equal(t, http.StatusNotFound, writer.Code)
		assert.Contains(t, writer.Body.String(), "failed to find live stream")
	})
}
