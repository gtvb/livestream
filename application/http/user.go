package http

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gtvb/livestream/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// Some Login description
// swagger:parameters loginUser
type LoginParamsWrapper struct {
	// in:body
	Body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
}

// Some Signup description
// swagger:parameters signupUser
type SignupParamsWrapper struct {
	// in:body
	Body struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
}

// swagger:response userResponse
type UserResponseWrapper struct {
	// in:body
	Body models.User
}

// swagger:response liveStreamsResponse
type LiveStreamsResponseWrapper struct {
	// in:body
	Body []models.LiveStream
}

// swagger:response tokenResponse
type TokenResponseWrapper struct {
	// O token JWT usado para próximas requisições protegidas
	Body struct {
		Token string `json:"token"`
	}
}

// swagger:response messageResponse
type MessageResponseWrapper struct {
	Body struct {
		Message string `json:"message"`
	}
}

// swagger:route POST /users/login users loginUser
// Realiza o login de um usuário e gera um token para futuras operações protegidas
//
// responses:
//
//	200: tokenResponse
//	400: messageResponse
//	404: messageResponse
//	500: messageResponse
func (env *ServerEnv) login(ctx *gin.Context) {
	var loginBody struct {
		Email    string
		Password string
	}

	if err := ctx.ShouldBindBodyWithJSON(&loginBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "failed to get request body"})
		return
	}

	// Check against database
	user, err := env.userRepository.GetUserByEmail(loginBody.Email)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			ctx.JSON(http.StatusNotFound, gin.H{"message": "a user with this email/password combination does not exist"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		}
		return
	}

	// Compare passwords
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginBody.Password))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "passwords don't match"})
		return
	}

	token, err := generateToken(user.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "failed to generate token"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": token})
}

// swagger:route POST /users/signup users signupUser
// Realiza o signup de um usuário e gera um token para futuras operações protegidas
// responses:
//
//	201: tokenResponse
//	400: messageResponse
//	500: messageResponse
func (env *ServerEnv) signup(ctx *gin.Context) {
	var signupBody struct {
		Name     string
		Email    string
		Password string
	}

	if err := ctx.ShouldBindBodyWithJSON(&signupBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "failed to get request body"})
		return
	}

	// Check against database
	user, err := env.userRepository.GetUserByEmail(signupBody.Email)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if user != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "a user with this email already exists"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(signupBody.Password), 14)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	newUser := models.NewUser(signupBody.Name, signupBody.Email, string(hashedPassword))
	userId, err := env.userRepository.CreateUser(newUser.Name, newUser.Email, newUser.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	token, err := generateToken(userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"token": token})
}

// swagger:route GET /users/:id users getUserProfile
// Retorna as informações sobre um usuário dado um id válido
//
// responses:
//
//	200: userResponse
//	404: messageResponse
func (env *ServerEnv) getUserProfile(ctx *gin.Context) {
	id := ctx.Param("id")
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	user, err := env.userRepository.GetUserById(objId)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "could not find a user with this id"})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

// swagger:route GET /livestream/:user_id users livestreams getUserLiveStreams
// Realiza o login de um usuário e gera um token para futuras operações protegidas
//
// responses:
//
//	200: liveStreamsResponse
//	404: messageResponse
func (env *ServerEnv) getUserLiveStreams(ctx *gin.Context) {
	id := ctx.Param("user_id")
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	livestreams, err := env.liveStreamsRepository.GetAllLiveStreamsByUserId(objId)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"livestreams": livestreams})
}

// swagger:route DELETE /user/:user_id users deleteUser
// Remove um usuário da base de dados, juntamente com todas as liveStreams cadastradas
// sobre seu mesmo id
//
// responses:
//
//	200: messageResponse
//	404: messageResponse
func (env *ServerEnv) deleteUser(ctx *gin.Context) {
	id := ctx.Param("user_id")
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

    err = env.liveStreamsRepository.DeleteLiveStreamsByPublisher(objId)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "failed to delete all streams for this user"})
		return
	}

	err = env.userRepository.DeleteUser(objId)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "failed to delete this user"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "user deleted"})
}

// swagger:route PATCH /user/:user_id users updateUser
// Atualiza os dados do usuário identificado pelo `user_id` especificado
// como parâmetro.
//
// responses:
//
//	200: messageResponse
//	400: messageResponse
func (env *ServerEnv) updateUser(ctx *gin.Context) {
	id := ctx.Param("user_id")
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	name := ctx.Query("name")
	if name == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "expected at least one query parameter for update"})
		return
	}

	err = env.userRepository.UpdateUserName(objId, name)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "falied to update this user's name"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "user updated"})
}

func (env *ServerEnv) getAllUsers(ctx *gin.Context) {
	users, err := env.userRepository.GetAllUsers()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "failed to fetch all users"})
		return
	}

	ctx.JSON(http.StatusOK, users)
}
