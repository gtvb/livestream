package http

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gtvb/livestream/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// swagger:route POST /livestreams/create livestreams createLiveStream
//
// Create a new live stream and assign it to the user specified in the request body.
//
// Responses:
//
//	201: liveStreamResponse
//	400: messageResponse
//	404: messageResponse
//	500: messageResponse
func (env *ServerEnv) createLiveStream(ctx *gin.Context) {
	ctx.Request.ParseMultipartForm(10 << 20)

	id := ctx.PostForm("publisher_id")
	name := ctx.PostForm("name")
	file, err := ctx.FormFile("thumbnail")

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "could not obtain the image"})
		return
	}

	userId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "invalid ID"})
		return
	}

	_, err = env.userRepository.GetUserById(userId)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "user not found"})
		return
	}

	streamKey, err := uuid.NewV7()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "failed to generate stream key"})
		return
	}

	filePath := filepath.Join("uploads", file.Filename)
	if err := ctx.SaveUploadedFile(file, filePath); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "failed to save the image"})
		return
	}

	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:3333"
	}

	fileUrl := fmt.Sprintf("%s/thumbs/%s", baseURL, file.Filename)
	streamId, err := env.liveStreamsRepository.CreateLiveStream(name, fileUrl, streamKey.String(), userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "failed to create the live stream"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"stream_id": streamId})
}

// swagger:route DELETE /livestreams/delete/{id} livestreams deleteLiveStream
//
// Delete a live stream given a valid `id`.
//
// Responses:
//
//	200: messageResponse
//	400: messageResponse
//	404: messageResponse
//	500: messageResponse
func (env *ServerEnv) deleteLiveStream(ctx *gin.Context) {
	streamId := ctx.Param("id")

	id, err := primitive.ObjectIDFromHex(streamId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "invalid ID"})
		return
	}

	err = env.liveStreamsRepository.DeleteLiveStream(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "failed to delete stream"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

// swagger:route PATCH /livestreams/update/{id} livestreams updateLiveStream
//
// Update the data of a live stream identified by the specified `id`.
//
// Responses:
//
//	200: messageResponse
//	400: messageResponse
//	404: messageResponse
//	500: messageResponse
func (env *ServerEnv) updateLiveStream(ctx *gin.Context) {
	id := ctx.Param("id")
	streamID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "unparseable ID"})
		return
	}

	var updateLiveStreamBody UpdateLiveStreamBody
	if err := ctx.ShouldBindBodyWithJSON(&updateLiveStreamBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "malformed request"})
		return
	}

	newData := bson.M{}
	if updateLiveStreamBody.Name != "" {
		newData["name"] = updateLiveStreamBody.Name
	}

	if updateLiveStreamBody.LiveStatus != nil {
		newData["live_stream_status"] = *updateLiveStreamBody.LiveStatus
	}

	err = env.liveStreamsRepository.UpdateLiveStream(streamID, newData)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "failed to update stream"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

// swagger:route GET /livestreams/info/{id} livestreams getLiveStreamData
//
// Get data for the live stream represented by the specified `id`.
//
// Responses:
//
//	200: liveStreamResponse
//	404: messageResponse
func (env *ServerEnv) getLiveStreamData(ctx *gin.Context) {
	streamId := ctx.Param("id")

	id, err := primitive.ObjectIDFromHex(streamId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "unparseable ID"})
		return
	}

	livestream, err := env.liveStreamsRepository.GetLiveStreamById(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "failed to find stream"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"livestream": livestream})
}

// generate swagger documentation for this function
// swagger:route GET /livestreams/feed livestreams getLiveStreamFeed
//
// Get a feed of live streams.
//
// Responses:
//
//	200: liveStreamFeedResponse
//	400: messageResponse
//	500: messageResponse
func (env *ServerEnv) getFeed(ctx *gin.Context) {
	numStreams := 20
	q := ctx.Query("q")
	if q != "" {
		qInt, err := strconv.Atoi(q)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "q needs to be an integer"})
			return
		}
		numStreams = qInt
	}

	livestreams, err := env.liveStreamsRepository.GetLiveStreamFeed(numStreams)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "failed to get livestream feed"})
		return
	}

	var users []*models.User
	for _, stream := range livestreams {
		user, err := env.userRepository.GetUserById(stream.PublisherId)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "could not get user for this stream"})
			return
		}

		users = append(users, user)
	}

	ctx.JSON(http.StatusOK, gin.H{"livestreams": livestreams, "users": users})
}

func (env *ServerEnv) getAllStreams(ctx *gin.Context) {
	livestreams, err := env.liveStreamsRepository.GetAllLiveStreams()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "failed to get livestreams"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"livestreams": livestreams})
}

func (env *ServerEnv) validateStream(ctx *gin.Context) {
	streamKey := ctx.Query("name")

	// Obter os parâmetros (gambiarra)
	tcurl := ctx.Query("swfurl")
	url, err := url.Parse(tcurl)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	username := url.Query().Get("username")
	password := url.Query().Get("password")

	if username == "" || password == "" {
		fmt.Println("No user and password")
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "missing username/password combination"})
		return
	}

	user, err := env.userRepository.GetUserByUsername(username)
	if err != nil {
		fmt.Println("No user")
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "invalid username"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		fmt.Println("Wrong password")
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "incorrect password"})
		return
	}

	ls, err := env.liveStreamsRepository.GetLiveStreamByStreamKey(streamKey)
	if err != nil {
		fmt.Println("No stream")
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "invalid stream key"})
		return
	}

	location := fmt.Sprintf("rtmp://127.0.0.1/hls-live/%s", ls.ID.Hex())

	ctx.Redirect(http.StatusFound, location)
}
