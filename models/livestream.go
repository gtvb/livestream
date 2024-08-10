package models

import "time"

type LiveStream struct {
	ID          int `bson:"_id"`
	Name        string
	StreamKey   string `bson:"stream_key"`
	ViewerCount int    `bson:"viewer_count"`

	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}

func NewLiveStream(name string) *LiveStream {
	return &LiveStream{
		Name:        name,
		ViewerCount: 0,

		CreatedAt: time.Now(),
	}
}
