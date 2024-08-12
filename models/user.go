package models

import (
	"time"
)

type UserRepositoryInterface interface {
	CreateUser(name, email, password string) (interface{}, error)
	DeleteUser(id int) (bool, error)
	UpdateUserName(id int, name string) (bool, error)
	UpdateUserEmail(id int, email string) (bool, error)
	UpdateUserPassword(id int, password string) (bool, error)
	UpdateUserAddLiveStream(id int, ls *LiveStream) (bool, error)
	GetUserByEmail(email string) (*User, error)
	GetUserById(id int) (*User, error)
	// Helper
	GetAllUsers() ([]*User, error)
}

type User struct {
	ID       int `bson:"_id"`
	Name     string
	Email    string
	Password string

	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`

	LiveStreams []*LiveStream `bson:"live_streams"`
}

func NewUser(name, email, password string) *User {
	return &User{
		Name:     name,
		Email:    email,
		Password: password,

		CreatedAt:   time.Now(),
		LiveStreams: make([]*LiveStream, 0),
	}
}

func (user *User) AddLiveStream(ls *LiveStream) {
	user.LiveStreams = append(user.LiveStreams, ls)
}
