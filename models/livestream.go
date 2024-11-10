package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LiveStreamRepositoryInterface interface {
	CreateLiveStream(name string, thumbnail string, streamKey string, publisherId primitive.ObjectID) (interface{}, error)
	DeleteLiveStream(id primitive.ObjectID) error
	DeleteLiveStreamsByPublisher(id primitive.ObjectID) error

	UpdateLiveStream(id primitive.ObjectID, newData bson.M) error

	IncrementLiveStreamUserCount(id primitive.ObjectID) error
	DecrementLiveStreamUserCount(id primitive.ObjectID) error

	GetAllLiveStreamsByUserId(id primitive.ObjectID) ([]*LiveStream, error)
	GetLiveStreamById(id primitive.ObjectID) (*LiveStream, error)
	GetLiveStreamByName(name string) (*LiveStream, error)
	GetLiveStreamByStreamKey(key string) (*LiveStream, error)
	GetLiveStreamFeed(maxStreams int) ([]*LiveStream, error)
	GetAllLiveStreams() ([]*LiveStream, error)
}

// Representa uma livestream acontecendo na plataforma
// swagger:model
type LiveStream struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name      string             `json:"name"`
	Thubmnail string             `json:"thumbnail"`

	StreamKey   string             `bson:"stream_key" json:"stream_key"`
	ViewerCount int                `bson:"viewer_count" json:"viewer_count"`
	PublisherId primitive.ObjectID `bson:"publisher_id" json:"publisher_id"`

	LiveStatus bool `bson:"live_stream_status" json:"live_stream_status"`

	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

func NewLiveStream(name string, thumbnail string, publisherId primitive.ObjectID, streamKey string) *LiveStream {
	return &LiveStream{
		Name:        name,
		Thubmnail:   thumbnail,
		PublisherId: publisherId,
		LiveStatus:  false,
		ViewerCount: 0,
		StreamKey:   streamKey,

		CreatedAt: time.Now(),
	}
}
