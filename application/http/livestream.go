package http

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
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

	_, err = env.liveStreamsRepository.CreateLiveStream(createLiveStreamBody.Name, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "livestream created"})
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
	streamId := ctx.Param("stream_id")

	id, err := primitive.ObjectIDFromHex(streamId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	err = env.liveStreamsRepository.DeleteLiveStream(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
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
	streamId := ctx.Param("stream_id")

	id, err := primitive.ObjectIDFromHex(streamId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	status := ctx.Query("status")
	name := ctx.Query("name")

	if status == "" && name == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "need one update paramater"})
		return
	}

	if status != "" {
		// TODO: verify error
		statusBool, _ := strconv.ParseBool(status)
		err := env.liveStreamsRepository.UpdateLiveStreamSetStatus(id, statusBool)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
	}

	if name != "" {
		// TODO: verify error
		err := env.liveStreamsRepository.UpdateLiveStreamName(id, name)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "succecss"})
}

func (env *ServerEnv) getLiveStreamData(ctx *gin.Context) {
	streamId := ctx.Param("stream_id")

	id, err := primitive.ObjectIDFromHex(streamId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	livestream, err := env.liveStreamsRepository.GetLiveStreamById(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"livestream": livestream})
}

func (env *ServerEnv) validateStream(ctx *gin.Context) {
	streamKey := ctx.Query("name")
	username := ctx.Query("username")
	password := ctx.Query("password")

	if username == "" || password == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "missing username/password combination"})
		return
	}

	// Verificar se usuário existe
	user, err := env.userRepository.GetUserByUsername(username)
	if err != nil {
		fmt.Println("No username")
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// Verificar se a senha está correta
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		fmt.Println("No password correct")
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "incorrect password"})
		return
	}

	// TODO: Utilizar a senha para criptografar a chave de stream. Após isso, verificar se
	// existe alguma entrada correspondente à esse hash. Se sim, a stream é válida

	fmt.Println(streamKey)
	ls, err := env.liveStreamsRepository.GetLiveStreamByStreamKey(streamKey)
	if err != nil {
		fmt.Println("No Key")
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// Redirecionar para a nova localização do arquivo
	location := fmt.Sprintf("rtmp://127.0.0.1/hls-live/%s", ls.ID.Hex())
	ctx.Redirect(http.StatusFound, location)
}
