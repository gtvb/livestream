package http

import (
	"net/http"
	"testing"

	"github.com/gtvb/livestream/models"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func createTestUser(env ServerEnv) *models.User {
	id, _ := env.userRepository.CreateUser("test", "test_username", "test@email.com", hashPassword("test_pass"))
	userID := id.(primitive.ObjectID)

	return &models.User{
		ID:       userID,
		Name:     "test",
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

	router := setupRouter(env)
	writer := makeRequest(router, "POST", "/livestreams/create", createLiveStreamBody, authParams)

	assert.Equal(t, http.StatusCreated, writer.Code)
	assert.Contains(t, writer.Body.String(), "livestream created")
}

func TestDeleteLiveStream(t *testing.T) {
	container := setupDatabase()
	defer container.Terminate()

	env := setupEnv(container.Database)
	user := createTestUser(env)

	authParams := &AuthParams{Email: user.Email, Password: user.Password}

	streamID, err := env.liveStreamsRepository.CreateLiveStream("Test Stream", user.ID)
	if err != nil {
		t.Fatalf("Failed to create test live stream: %v", err)
	}
	id := streamID.(primitive.ObjectID)

	router := setupRouter(env)
	writer := makeRequest(router, "DELETE", "/livestreams/delete/"+id.Hex(), nil, authParams)

	assert.Equal(t, http.StatusOK, writer.Code)
	assert.Contains(t, writer.Body.String(), "success")
}

func TestUpdateLiveStream(t *testing.T) {
	container := setupDatabase()
	defer container.Terminate()

	env := setupEnv(container.Database)
	user := createTestUser(env)

	authParams := &AuthParams{Email: user.Email, Password: user.Password}

	streamID, err := env.liveStreamsRepository.CreateLiveStream("Old Name", user.ID)
	if err != nil {
		t.Fatalf("Failed to create test live stream: %v", err)
	}
	id := streamID.(primitive.ObjectID)

	router := setupRouter(env)
	writer := makeRequest(router, "PATCH", "/livestreams/update/"+id.Hex()+"?name=New Name", nil, authParams)

	assert.Equal(t, http.StatusOK, writer.Code)
	assert.Contains(t, writer.Body.String(), "success")
}

func TestGetLiveStreamInfo(t *testing.T) {
	container := setupDatabase()
	defer container.Terminate()

	env := setupEnv(container.Database)
	user := createTestUser(env)

	authParams := &AuthParams{Email: user.Email, Password: user.Password}

	streamID, err := env.liveStreamsRepository.CreateLiveStream("Test Stream", user.ID)
	if err != nil {
		t.Fatalf("Failed to create test live stream: %v", err)
	}

	id := streamID.(primitive.ObjectID)

	router := setupRouter(env)
	writer := makeRequest(router, "GET", "/livestreams/info/"+id.Hex(), nil, authParams)

	assert.Equal(t, http.StatusOK, writer.Code)
	assert.Contains(t, writer.Body.String(), "Test Stream")
}
