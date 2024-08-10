package models

import "time"

type User struct {
	ID       int `bson:"_id"`
	Name     string
	Email    string
	Password string

	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`

	LiveStreams []LiveStream `bson:"live_streams"`
}

func NewUser(name, email, password string) *User {
	return &User{
		Name:     name,
		Email:    email,
		Password: password,

		CreatedAt:   time.Now(),
		LiveStreams: make([]LiveStream, 0),
	}
}
