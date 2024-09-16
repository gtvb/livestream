package http

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type CreateLiveStreamBody struct {
	// User ID of the stream creator
	// required: true
	UserId string `json:"user_id"`
	// Name of the live stream
	// required: true
	Name string `json:"name"`
}

// CreateLiveStreamParamsWrapper contains parameters for creating a live stream.
// swagger:parameters createLiveStream
type CreateLiveStreamParamsWrapper struct {
	// in:body
	Body CreateLiveStreamBody
}

// LiveStreamResponseWrapper contains a response with live stream data.
// swagger:response liveStreamResponse
type LiveStreamResponseWrapper struct {
	// in:body
	Body struct {
		// ID of the live stream
		StreamId primitive.ObjectID `json:"stream_id"`
	}
}

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
	var createLiveStreamBody CreateLiveStreamBody

	if err := ctx.ShouldBindBodyWithJSON(&createLiveStreamBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "could not get body from request"})
		return
	}

	id, err := primitive.ObjectIDFromHex(createLiveStreamBody.UserId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	_, err = env.userRepository.GetUserById(id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			ctx.JSON(http.StatusNotFound, gin.H{"message": "a user with this id does not exist"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		}
		return
	}

	_, err = env.liveStreamsRepository.CreateLiveStream(createLiveStreamBody.Name, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "livestream created"})
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
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	err = env.liveStreamsRepository.DeleteLiveStream(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "failed to delete live stream"})
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
	streamId := ctx.Param("id")

	id, err := primitive.ObjectIDFromHex(streamId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	status := ctx.Query("status")
	name := ctx.Query("name")

	if status == "" && name == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "need one update parameter"})
		return
	}

	if status != "" {
		statusBool, err := strconv.ParseBool(status)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "invalid status value"})
			return
		}

		err = env.liveStreamsRepository.UpdateLiveStreamSetStatus(id, statusBool)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
	}

	if name != "" {
		err = env.liveStreamsRepository.UpdateLiveStreamName(id, name)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
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
//	500: messageResponse
func (env *ServerEnv) getLiveStreamData(ctx *gin.Context) {
	streamId := ctx.Param("id")

	id, err := primitive.ObjectIDFromHex(streamId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	livestream, err := env.liveStreamsRepository.GetLiveStreamById(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "failed to find live stream"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"livestream": livestream})
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
