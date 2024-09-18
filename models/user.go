package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserRepositoryInterface interface {
	CreateUser(username, email, password string) (interface{}, error)
	DeleteUser(id primitive.ObjectID) error

	UpdateUser(id primitive.ObjectID, newData bson.M) error

	UpdateUserAddToFollowList(id primitive.ObjectID, following primitive.ObjectID) error
	UpdateUserRemoveFromFollowList(id primitive.ObjectID, following primitive.ObjectID) error

	GetUserById(id primitive.ObjectID) (*User, error)
	GetUserByEmail(email string) (*User, error)
	GetUserByUsername(username string) (*User, error)

	GetAllUsers() ([]*User, error)
}

// Representa um usuário cadastrado na plataforma
// swagger:model
type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username string             `bson:"username" json:"username"`
	Email    string             `bson:"email" json:"email"`
	Password string             `bson:"password" json:"password"`

	Following []primitive.ObjectID `bson:"following" json:"following"`

	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

func NewUser(username, email, password string) *User {
	return &User{
		Username: username,
		Email:    email,
		Password: password,

		Following: make([]primitive.ObjectID, 0),

		CreatedAt: time.Now(),
	}
}
