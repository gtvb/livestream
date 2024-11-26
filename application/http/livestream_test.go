package http

import (
	"net/http"
	"strconv"
	"testing"

	"github.com/gtvb/livestream/models"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
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

// func TestCreateLiveStream(t *testing.T) {
// 	container := setupDatabase()
// 	defer container.Terminate()

// 	env := setupEnv(container.Database)
// 	router := setupRouter(env)
// 	user := createTestUser(env)

// 	createLiveStreamBody := CreateLiveStreamBody{
// 		UserId: user.ID.Hex(),
// 		Name:   "Test Live Stream",
// 	}

// 	t.Run("Correct body", func(t *testing.T) {
// 		writer := makeRequest(router, "POST", "/livestreams/create", createLiveStreamBody)
// 		assert.Equal(t, http.StatusCreated, writer.Code)
// 		assert.Contains(t, writer.Body.String(), "stream_id")
// 	})

// 	t.Run("Invalid ID", func(t *testing.T) {
// 		createLiveStreamBody.UserId = "invalidID"
// 		writer := makeRequest(router, "POST", "/livestreams/create", createLiveStreamBody)
// 		assert.Equal(t, http.StatusBadRequest, writer.Code)
// 		assert.Contains(t, writer.Body.String(), "invalid ObjectID")
// 	})

// 	t.Run("No body on request", func(t *testing.T) {
// 		writer := makeRequest(router, "POST", "/livestreams/create", nil)
// 		assert.Equal(t, http.StatusBadRequest, writer.Code)
// 		assert.Contains(t, writer.Body.String(), "could not get body from request")
// 	})
// }

func TestDeleteLiveStream(t *testing.T) {
	container := setupDatabase()
	defer container.Terminate()

	env := setupEnv(container.Database)
	router := setupRouter(env)

	user := createTestUser(env)

	t.Run("Successfully create stream", func(t *testing.T) {
		streamID, err := env.liveStreamsRepository.CreateLiveStream("Test Stream", "fake-thumbnail", "streamkey-test", user.ID)
		if err != nil {
			t.Fatalf("Failed to create test live stream: %v", err)
		}
		id := streamID.(primitive.ObjectID)

		writer := makeRequest(router, "DELETE", "/livestreams/delete/"+id.Hex(), nil)
		assert.Equal(t, http.StatusOK, writer.Code)
		assert.Contains(t, writer.Body.String(), "success")
	})

	t.Run("Invalid ID", func(t *testing.T) {
		writer := makeRequest(router, "DELETE", "/livestreams/delete/invalidID", nil)
		assert.Equal(t, http.StatusBadRequest, writer.Code)
		assert.Contains(t, writer.Body.String(), "invalid ID")
	})

	t.Run("Stream not found", func(t *testing.T) {
		writer := makeRequest(router, "DELETE", "/livestreams/delete/"+primitive.NewObjectID().Hex(), nil)
		assert.Equal(t, http.StatusNotFound, writer.Code)
		assert.Contains(t, writer.Body.String(), "failed to delete stream")
	})
}

func TestUpdateLiveStream(t *testing.T) {
	container := setupDatabase()
	defer container.Terminate()

	env := setupEnv(container.Database)
	router := setupRouter(env)
	user := createTestUser(env)

	streamID, _ := env.liveStreamsRepository.CreateLiveStream("Test Stream", "fake-thumbnail", "streamkey-test", user.ID)
	id := streamID.(primitive.ObjectID)

	t.Run("Correct body", func(t *testing.T) {
		updateLiveStreamBody := UpdateLiveStreamBody{
			Name: "Updated Stream Name",
		}

		writer := makeRequest(router, "PATCH", "/livestreams/update/"+id.Hex(), updateLiveStreamBody)
		assert.Equal(t, http.StatusOK, writer.Code)
		assert.Contains(t, writer.Body.String(), "success")
	})

	t.Run("Invalid ID", func(t *testing.T) {
		writer := makeRequest(router, "PATCH", "/livestreams/update/invalidID", nil)
		assert.Equal(t, http.StatusBadRequest, writer.Code)
		assert.Contains(t, writer.Body.String(), "unparseable ID")
	})

	t.Run("Unexistant ID", func(t *testing.T) {
		writer := makeRequest(router, "PATCH", "/livestreams/update/"+primitive.NewObjectID().Hex(), nil)
		assert.Equal(t, http.StatusBadRequest, writer.Code)
		assert.Contains(t, writer.Body.String(), "failed to update stream")
	})
}

func TestGetLiveStreamData(t *testing.T) {
	container := setupDatabase()
	defer container.Terminate()

	env := setupEnv(container.Database)
	router := setupRouter(env)

	user := createTestUser(env)
	streamID, _ := env.liveStreamsRepository.CreateLiveStream("Test Stream", "fake-thumbnail", "streamkey-test", user.ID)

	t.Run("Correct ID", func(t *testing.T) {
		writer := makeRequest(router, "GET", "/livestreams/info/"+streamID.(primitive.ObjectID).Hex(), nil)
		assert.Equal(t, http.StatusOK, writer.Code)
		assert.Contains(t, writer.Body.String(), "Test Stream")
	})

	t.Run("Invalid ID", func(t *testing.T) {
		writer := makeRequest(router, "GET", "/livestreams/info/invalidID", nil)
		assert.Equal(t, http.StatusBadRequest, writer.Code)
		assert.Contains(t, writer.Body.String(), "unparseable ID")
	})

	t.Run("Unexistant ID", func(t *testing.T) {
		writer := makeRequest(router, "GET", "/livestreams/info/"+primitive.NewObjectID().Hex(), nil)
		assert.Equal(t, http.StatusNotFound, writer.Code)
		assert.Contains(t, writer.Body.String(), "failed to find stream")
	})
}

func TestGetLiveStreamFeed(t *testing.T) {
	container := setupDatabase()
	defer container.Terminate()

	env := setupEnv(container.Database)
	router := setupRouter(env)
	user := createTestUser(env)

	// Create some live streams
	for i := 0; i < 5; i++ {
		streamID, _ := env.liveStreamsRepository.CreateLiveStream("Test Stream "+strconv.Itoa(i), "fake-thumbnail", "streamkey-test"+strconv.Itoa(i), user.ID)
		id := streamID.(primitive.ObjectID)

		newData := bson.M{"live_stream_status": true}
		env.liveStreamsRepository.UpdateLiveStream(id, newData)
	}

	t.Run("DefaultNumStreams", func(t *testing.T) {
		writer := makeRequest(router, "GET", "/livestreams/feed", nil)
		assert.Equal(t, http.StatusOK, writer.Code)
		assert.Contains(t, writer.Body.String(), "Test Stream")
	})

	t.Run("CustomNumStreams", func(t *testing.T) {
		writer := makeRequest(router, "GET", "/livestreams/feed?q=3", nil)
		assert.Equal(t, http.StatusOK, writer.Code)
		assert.Contains(t, writer.Body.String(), "Test Stream")
	})

	t.Run("InvalidQueryParam", func(t *testing.T) {
		writer := makeRequest(router, "GET", "/livestreams/feed?q=invalid", nil)
		assert.Equal(t, http.StatusBadRequest, writer.Code)
		assert.Contains(t, writer.Body.String(), "q needs to be an integer")
	})
}
