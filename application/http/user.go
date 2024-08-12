package http

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/gtvb/livestream/models"
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

// swagger:response tokenResponse
type TokenResponseWrapper struct {
	// O token JWT usado para próximas requisições protegidas
	Body struct {
		Token string `json:"token"`
	}
}

// swagger:route POST /users/login users loginUser
// Realiza o login de um usuário e gera um token para futuras operações protegidas
//
// responses:
//
//	200: tokenResponse
func (env *ServerEnv) login(w http.ResponseWriter, req *http.Request) {
	var loginBody struct {
		Email    string
		Password string
	}

	err := json.NewDecoder(req.Body).Decode(&loginBody)
	if err != nil {
		log.Printf("Decode: %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check against database
	user, err := env.userRepository.GetUserByEmail(loginBody.Email)
	if err != nil {
		log.Printf("GetUserByEmail: %s\n", err.Error())
		if errors.Is(err, mongo.ErrNoDocuments) {
			http.Error(w, "could not find a user with this email/password", http.StatusNotFound)
		} else {
			http.Error(w, "error while fetching user", http.StatusInternalServerError)
		}
		return
	}

	// Compare passwords
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginBody.Password))
	if err != nil {
		log.Printf("bcrypt: %s\n", err.Error())
		http.Error(w, "wrong email or password", http.StatusNotFound)
		return
	}

	token, err := generateToken(user.ID)
	if err != nil {
		log.Printf("token: %s\n", err.Error())
		http.Error(w, "failed to register user", http.StatusInternalServerError)
	}

	w.Write([]byte(token))
}

// swagger:route POST /users/signup users signupUser
// Realiza o signup de um usuário e gera um token para futuras operações protegidas
// responses:
//
//	200: tokenResponse
func (env *ServerEnv) signup(w http.ResponseWriter, req *http.Request) {
	var signupBody struct {
		Name     string
		Email    string
		Password string
	}

	err := json.NewDecoder(req.Body).Decode(&signupBody)
	if err != nil {
		log.Printf("Decode: %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Check against database
	user, err := env.userRepository.GetUserByEmail(signupBody.Email)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		log.Printf("1 GetUserByEmail: %s\n", err.Error())
		http.Error(w, "failed to process request", http.StatusInternalServerError)
		return
	}

	if user != nil {
		log.Printf("no user: %s\n", err.Error())
		http.Error(w, "a user with this email already exists", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(signupBody.Password), 14)
	if err != nil {
		log.Printf("bcrypt: %s\n", err.Error())
		http.Error(w, "failed to process request", http.StatusInternalServerError)
		return
	}

	newUser := models.NewUser(signupBody.Name, signupBody.Email, string(hashedPassword))
	_, err = env.userRepository.CreateUser(newUser.Name, newUser.Email, newUser.Password)
	if err != nil {
		log.Printf("CreateUser: %s\n", err.Error())
		http.Error(w, "failed to register user", http.StatusInternalServerError)
		return
	}

	token, err := generateToken(newUser.ID)
	if err != nil {
		log.Printf("token: %s\n", err.Error())
		http.Error(w, "failed to register user", http.StatusInternalServerError)
	}

	w.Write([]byte(token))
}

// func (env *ServerEnv) allUsers(w http.ResponseWriter, req *http.Request) {
// 	w.Header().Add("Content-Type", "application/json")
// 	users, err := env.userRepository.GetAllUsers()
// 	if err != nil {
// 		log.Print(err.Error())
// 		http.Error(w, "could not get all users", http.StatusInternalServerError)
// 		return
// 	}

// 	usersJson, err := json.Marshal(users)
// 	if err != nil {
// 		log.Print(err.Error())
// 		http.Error(w, "could not JSON encode all users", http.StatusInternalServerError)
// 		return
// 	}

// 	w.Write(usersJson)
// }
