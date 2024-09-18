package http

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gtvb/livestream/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type UpdateUserBody struct {
	// User's username
	// required: true
	Username string `json:"username"`
	// User's email
	// required: true
	Email string `json:"email"`
	// User's password
	// required: true
	Password string `json:"password"`
}

// UpdateUserParamsWrapper contains parameters for updating a user
// swagger:parameters updateUser
type UpdateUserParamsWrapper struct {
	// in:body
	Body UpdateUserBody
}

type LoginBody struct {
	// User's email
	// required: true
	Email string `json:"email"`
	// User's password
	// required: true
	Password string `json:"password"`
}

// LoginParamsWrapper contains parameters for user login.
// swagger:parameters loginUser
type LoginParamsWrapper struct {
	// in:body
	Body LoginBody
}

type SignupBody struct {
	// User's name
	// required: true
	Name string `json:"name"`
	// User's username
	// required: true
	Username string `json:"username"`
	// User's email
	// required: true
	Email string `json:"email"`
	// User's password
	// required: true
	Password string `json:"password"`
}

// SignupParamsWrapper contains parameters for user signup.
// swagger:parameters signupUser
type SignupParamsWrapper struct {
	// in:body
	Body SignupBody
}

// UserResponseWrapper contains a user response.
// swagger:response userResponse
type UserResponseWrapper struct {
	// in:body
	Body struct {
		// The user details
		User models.User `json:"user"`
	}
}

// UserListResponseWrapper contains a user list response.
// swagger:response userListResponse
type UserListResponseWrapper struct {
	// in:body
	Body struct {
		// The user details
		Users []models.User `json:"users"`
	}
}

// LiveStreamsResponseWrapper contains a response with live streams.
// swagger:response liveStreamsResponse
type LiveStreamsResponseWrapper struct {
	// in:body
	Body struct {
		// List of live streams
		LiveStreams []models.LiveStream `json:"livestreams"`
	}
}

// TokenResponseWrapper contains a token response.
// swagger:response tokenResponse
type TokenResponseWrapper struct {
	// The JWT token for future protected requests.
	// required: true
	Body struct {
		Token string `json:"token"`
	}
}

// MessageResponseWrapper contains a message response.
// swagger:response messageResponse
type MessageResponseWrapper struct {
	Body struct {
		// A descriptive message
		Message string `json:"message"`
	}
}

// swagger:route POST /users/login users loginUser
//
// Login a user and generate a token for future protected operations.
//
// Responses:
//
//	200: tokenResponse
//	400: messageResponse
//	404: messageResponse
//	500: messageResponse
func (env *ServerEnv) login(ctx *gin.Context) {
	var loginBody LoginBody

	if err := ctx.ShouldBindBodyWithJSON(&loginBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "failed to get request body"})
		return
	}

	user, err := env.userRepository.GetUserByEmail(loginBody.Email)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			ctx.JSON(http.StatusNotFound, gin.H{"message": "a user with this email/password combination does not exist"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		}
		return
	}

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
//
// Signup a user and generate a token for future protected operations.
//
// Responses:
//
//	201: tokenResponse
//	400: messageResponse
//	500: messageResponse
func (env *ServerEnv) signup(ctx *gin.Context) {
	var signupBody SignupBody

	if err := ctx.ShouldBindBodyWithJSON(&signupBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "failed to get request body"})
		return
	}

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

	userId, err := env.userRepository.CreateUser(signupBody.Username, signupBody.Email, string(hashedPassword))
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

// swagger:route GET /users/{id} users getUserProfile
//
// Get user profile information given a valid id.
//
// Responses:
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

	ctx.JSON(http.StatusOK, gin.H{"user": user})
}

// swagger:route GET /livestreams/{user_id} livestreams getUserLiveStreams
//
// Get all live streams that belong to the user specified by `user_id`.
//
// Responses:
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

// swagger:route DELETE /users/{id} users deleteUser
//
// Delete a user from the database along with all their registered live streams.
//
// Responses:
//
//	200: messageResponse
//	404: messageResponse
func (env *ServerEnv) deleteUser(ctx *gin.Context) {
	id := ctx.Param("id")
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

// swagger:route PATCH /users/{id} users updateUser
//
// Update the user's data identified by the specified `id` parameter.
//
// Responses:
//
//	200: messageResponse
//	400: messageResponse
func (env *ServerEnv) updateUser(ctx *gin.Context) {
	id := ctx.Param("id")
	userID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	var updateBody UpdateUserBody
	if err := ctx.ShouldBindBodyWithJSON(&updateBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "could not parse request body: " + err.Error()})
		return
	}

	newData := bson.M{}
	if updateBody.Email != "" {
		newData["email"] = updateBody.Email
	}

	if updateBody.Username != "" {
		newData["username"] = updateBody.Username
	}

	if updateBody.Password != "" {
		newData["password"] = updateBody.Email
	}

	err = env.userRepository.UpdateUser(userID, newData)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "could not update user: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "updated user with success"})
}

// swagger:route GET /users/all users getAllUsers
//
// Get all users.
//
// Responses:
//
//	200: []User
//	500: messageResponse
func (env *ServerEnv) getAllUsers(ctx *gin.Context) {
	users, err := env.userRepository.GetAllUsers()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "failed to fetch all users"})
		return
	}

	ctx.JSON(http.StatusOK, users)
}
