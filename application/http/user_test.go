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
		Name:     "test",
		Username: "test_username",
		Password: "test_pass",
	}

	router := setupRouter(env)
	writer := makeRequest(router, "POST", "/user/signup", signupBody, nil)

	assert.Equal(t, http.StatusCreated, writer.Code)
	assert.Contains(t, writer.Body.String(), "token")
}

func TestUserLogin(t *testing.T) {
	container := setupDatabase()
	defer container.Terminate()

	env := setupEnv(container.Database)
	env.userRepository.CreateUser("test", "test_username", "test@email.com", hashPassword("test_pass"))

	loginBody := LoginBody{
		Email:    "test@email.com",
		Password: "test_pass",
	}

	router := setupRouter(env)
	writer := makeRequest(router, "POST", "/user/login", loginBody, nil)

	assert.Equal(t, http.StatusOK, writer.Code)
	assert.Contains(t, writer.Body.String(), "token")
}

func TestGetUserProfile(t *testing.T) {
	container := setupDatabase()
	defer container.Terminate()

	env := setupEnv(container.Database)

	id, _ := env.userRepository.CreateUser("test", "test_username", "test@email.com", hashPassword("test_pass"))
	userID := id.(primitive.ObjectID)

	router := setupRouter(env)
	writer := makeRequest(router, "GET", "/user/"+userID.Hex(), nil, &AuthParams{Email: "test@email.com", Password: "test_pass"})

	assert.Equal(t, http.StatusOK, writer.Code)
	assert.Contains(t, writer.Body.String(), "test")
}

func TestDeleteUser(t *testing.T) {
	container := setupDatabase()
	defer container.Terminate()

	env := setupEnv(container.Database)

	id, _ := env.userRepository.CreateUser("test", "test_username", "test@email.com", hashPassword("test_pass"))
	userID := id.(primitive.ObjectID)

	router := setupRouter(env)
	writer := makeRequest(router, "DELETE", "/user/delete/"+userID.Hex(), nil, &AuthParams{Email: "test@email.com", Password: "test_pass"})

	assert.Equal(t, http.StatusOK, writer.Code)
	assert.Contains(t, writer.Body.String(), "user deleted")
}

func TestGetAllUsers(t *testing.T) {
	container := setupDatabase()
	defer container.Terminate()

	env := setupEnv(container.Database)

	users := []models.User{
		{Name: "test1", Username: "test_username1", Email: "test1@email.com", Password: hashPassword("test_pass1")},
		{Name: "test2", Username: "test_username2", Email: "test2@email.com", Password: hashPassword("test_pass2")},
		{Name: "test3", Username: "test_username3", Email: "test3@email.com", Password: hashPassword("test_pass3")},
	}

	for _, user := range users {
		env.userRepository.CreateUser(user.Name, user.Username, user.Email, user.Password)
	}

	router := setupRouter(env)
	writer := makeRequest(router, "GET", "/user/all", nil, &AuthParams{Email: "test@email.com", Password: "test_pass"})

	assert.Equal(t, http.StatusOK, writer.Code)
	// assert.JSONEq(t, string(expectedJSON), writer.Body.String())
}

func hashPassword(password string) string {
	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashed)
}
