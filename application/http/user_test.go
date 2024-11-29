package http

import (
	"net/http"
	"testing"

	"github.com/gtvb/livestream/models"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func TestUserSignup(t *testing.T) {
	container := setupDatabase()
	defer container.Terminate()

	env := setupEnv(container.Database)

	signupBody := SignupBody{
		Email:    "test@email.com",
		Username: "test_username",
		Password: "test_pass",
	}

	router := setupRouter(env)
	writer := makeRequest(router, "POST", "/user/signup", signupBody)

	assert.Equal(t, http.StatusCreated, writer.Code)
	assert.Contains(t, writer.Body.String(), "success")
}

func TestUserLogin(t *testing.T) {
	container := setupDatabase()
	defer container.Terminate()

	env := setupEnv(container.Database)
	env.userRepository.CreateUser("test_username", "test@email.com", hashPassword("test_pass"))

	loginBody := LoginBody{
		Email:    "test@email.com",
		Password: "test_pass",
	}

	router := setupRouter(env)
	writer := makeRequest(router, "POST", "/user/login", loginBody)

	assert.Equal(t, http.StatusOK, writer.Code)
	assert.Contains(t, writer.Body.String(), "test_username")
}

func TestGetUserProfile(t *testing.T) {
	container := setupDatabase()
	defer container.Terminate()

	env := setupEnv(container.Database)

	id, _ := env.userRepository.CreateUser("test_username", "test@email.com", hashPassword("test_pass"))
	userID := id.(primitive.ObjectID)

	router := setupRouter(env)
	writer := makeRequest(router, "GET", "/user/"+userID.Hex(), nil)

	assert.Equal(t, http.StatusOK, writer.Code)
	assert.Contains(t, writer.Body.String(), "test")
}

func TestDeleteUser(t *testing.T) {
	container := setupDatabase()
	defer container.Terminate()

	env := setupEnv(container.Database)

	id, _ := env.userRepository.CreateUser("test_username", "test@email.com", hashPassword("test_pass"))
	userID := id.(primitive.ObjectID)

	router := setupRouter(env)
	writer := makeRequest(router, "DELETE", "/user/delete/"+userID.Hex(), nil)

	assert.Equal(t, http.StatusOK, writer.Code)
	assert.Contains(t, writer.Body.String(), "success")
}

func TestUpdateUser(t *testing.T) {
	container := setupDatabase()
	defer container.Terminate()

	env := setupEnv(container.Database)

	id, _ := env.userRepository.CreateUser("test_username", "test@email.com", hashPassword("test_pass"))
	userID := id.(primitive.ObjectID)

	updateBody := UpdateUserBody{
		Username: "test_username_new",
		Email:    "test@new.email.com",
	}

	router := setupRouter(env)
	writer := makeRequest(router, "PATCH", "/user/update/"+userID.Hex(), updateBody)

	assert.Equal(t, http.StatusOK, writer.Code)
	assert.Contains(t, writer.Body.String(), "success")
}

func TestFollowUser(t *testing.T) {
	container := setupDatabase()
	defer container.Terminate()

	env := setupEnv(container.Database)

	id1, _ := env.userRepository.CreateUser("test_username1", "test1@email.com", hashPassword(("test1")))
	id2, _ := env.userRepository.CreateUser("test_username2", "test2@email.com", hashPassword(("test2")))
	user1ID := id1.(primitive.ObjectID)
	user2ID := id2.(primitive.ObjectID)

	followBody := FollowBody{
		UserID: user1ID,
	}

	router := setupRouter(env)
	writer := makeRequest(router, "PATCH", "/user/follow/"+user2ID.Hex(), followBody)

	assert.Equal(t, http.StatusOK, writer.Code)
	assert.Contains(t, writer.Body.String(), "success")
}

func TestUnFollowUser(t *testing.T) {
	container := setupDatabase()
	defer container.Terminate()

	env := setupEnv(container.Database)

	id1, _ := env.userRepository.CreateUser("test_username1", "test1@email.com", hashPassword(("test1")))
	id2, _ := env.userRepository.CreateUser("test_username2", "test2@email.com", hashPassword(("test2")))
	user1ID := id1.(primitive.ObjectID)
	user2ID := id2.(primitive.ObjectID)

	followBody := FollowBody{
		UserID: user1ID,
	}

	router := setupRouter(env)
	writer := makeRequest(router, "PATCH", "/user/unfollow/"+user2ID.Hex(), followBody)

	assert.Equal(t, http.StatusOK, writer.Code)
	assert.Contains(t, writer.Body.String(), "success")
}

func TestGetAllUsers(t *testing.T) {
	container := setupDatabase()
	defer container.Terminate()

	env := setupEnv(container.Database)

	users := []models.User{
		{Username: "test_username1", Email: "test1@email.com", Password: hashPassword("test_pass1")},
		{Username: "test_username2", Email: "test2@email.com", Password: hashPassword("test_pass2")},
		{Username: "test_username3", Email: "test3@email.com", Password: hashPassword("test_pass3")},
	}

	for _, user := range users {
		env.userRepository.CreateUser(user.Username, user.Email, user.Password)
	}

	router := setupRouter(env)
	writer := makeRequest(router, "GET", "/user/all", nil)

	assert.Equal(t, http.StatusOK, writer.Code)
	// assert.JSONEq(t, string(expectedJSON), writer.Body.String())
}

func hashPassword(password string) string {
	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashed)
}
