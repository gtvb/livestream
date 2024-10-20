package http

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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

	// Get the user responsibe foe this livestream
	_, err = env.userRepository.GetUserById(id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			ctx.JSON(http.StatusNotFound, gin.H{"message": "a user with this id does not exist"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		}
		return
	}

	// Generate a streamkey for this user, based on his password
	streamKey, err := uuid.NewV7()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// // Verify the password
	// if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(createLiveStreamBody.Password)); err != nil {
	// 	ctx.JSON(http.StatusBadRequest, gin.H{"message": "incorrect password"})
	// 	return
	// }

	// // TODO: Encrypt the streamkey using the password

	_, err = env.liveStreamsRepository.CreateLiveStream(createLiveStreamBody.Name, streamKey.String(), id)
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
	id := ctx.Param("id")
	streamID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	var updateLiveStreamBody UpdateLiveStreamBody
	if err := ctx.ShouldBindBodyWithJSON(&updateLiveStreamBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "não foi possível interpretar o corpo da requisição"})
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
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "falha ao atualizar a live"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "live atualizada com sucesso"})
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
