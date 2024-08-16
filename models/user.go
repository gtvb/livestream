package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserRepositoryInterface interface {
	CreateUser(name, email, password string) (interface{}, error)
	DeleteUser(id primitive.ObjectID) error

	UpdateUserName(id primitive.ObjectID, name string) error
	UpdateUserEmail(id primitive.ObjectID, email string) error
	UpdateUserPassword(id primitive.ObjectID, password string) error

	GetUserByEmail(email string) (*User, error)
	GetUserById(id primitive.ObjectID) (*User, error)
	GetAllUsers() ([]*User, error)
}

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name     string             `json:"name"`
	Email    string             `json:"email"`
	Password string             `json:"password"`

	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

func NewUser(name, email, password string) *User {
	return &User{
		Name:     name,
		Email:    email,
		Password: password,

		CreatedAt: time.Now(),
	}
}
