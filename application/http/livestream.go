package http

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gtvb/livestream/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// swagger:parameters createLiveStream
type CreateLiveStreamParamsWrapper struct {
	// in:body
	Body struct {
		UserId string `json:"user_id"`
		Name   string `json:"name"`
	}
}

// swagger:response liveStreamResponse
type LiveStreamResponseWrapper struct {
	// in:body
	Body struct {
		StreamId primitive.ObjectID `json:"stream_id"`
	}
}

// swagger:route POST /livestream/create livestreams createLiveStream
//
// Cria uma nova livestream e atribui ela ao usuário contido no
// corpo da requisição.
// responses:
// 200: liveStreamResponse
// 404: messageResponse
// 500: messageResponse
func (env *ServerEnv) createLiveStream(ctx *gin.Context) {
	var createLiveStreamBody struct {
		UserId string `json:"user_id"`
		Name   string `json:"name"`
	}

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

	liveStream := models.NewLiveStream(createLiveStreamBody.Name, id)
	_, err = env.liveStreamsRepository.CreateLiveStream(liveStream.Name, liveStream.PublisherId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"livestream_data": liveStream})
}

// swagger:route DELETE /livestream/delete/:id livestreams deleteLiveStream
//
// Deleta uma livestream dado um id válido.
//
// responses:
// 200: liveStreamResponse
// 404: messageResponse
// 500: messageResponse
func (env *ServerEnv) deleteLiveStream(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "test"})
}

// swagger:route PATCH /livestream/update/:id livestreams updateLiveStream
//
// Cria uma nova livestream e atribui ela ao usuário contido no
// corpo da requisição.
// responses:
// 200: liveStreamResponse
// 404: messageResponse
// 500: messageResponse
func (env *ServerEnv) updateLiveStream(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "test"})
}
