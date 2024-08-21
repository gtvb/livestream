package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserRepositoryInterface interface {
	CreateUser(name, username, email, password string) (interface{}, error)
	DeleteUser(id primitive.ObjectID) error

	UpdateUserName(id primitive.ObjectID, name string) error
	UpdateUserEmail(id primitive.ObjectID, email string) error
	UpdateUserPassword(id primitive.ObjectID, password string) error

	GetUserById(id primitive.ObjectID) (*User, error)
	GetUserByEmail(email string) (*User, error)
	GetUserByUsername(username string) (*User, error)

	GetAllUsers() ([]*User, error)
}

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name     string             `bson:"name" json:"name"`
	Username string             `bson:"username" json:"username"`
	Email    string             `bson:"email" json:"email"`
	Password string             `bson:"password" json:"password"`

	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

func NewUser(name, username, email, password string) *User {
	return &User{
		Name:     name,
		Username: username,
		Email:    email,
		Password: password,

		CreatedAt: time.Now(),
	}
}
