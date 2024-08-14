package models

import (
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LiveStreamRepositoryInterface interface {
	CreateLiveStream(name string, publisherId primitive.ObjectID) (interface{}, error)
	DeleteLiveStream(id primitive.ObjectID) (bool, error)
	DeleteLiveStreamsByPublisher(id primitive.ObjectID) (bool, error)

	UpdateLiveStreamName(id primitive.ObjectID, name string) (bool, error)
	IncrementLiveStreamUserCount(id primitive.ObjectID) (bool, error)
	DecrementLiveStreamUserCount(id primitive.ObjectID) (bool, error)

	GetAllLiveStreamsByUserId(id primitive.ObjectID) ([]*LiveStream, error)
	GetLiveStreamByName(name string) (*LiveStream, error)
	GetAllLiveStreams() ([]*LiveStream, error)
}

type LiveStream struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string             `json:"name"`
	StreamKey   string             `bson:"stream_key" json:"stream_key"`
	ViewerCount int                `bson:"viewer_count" json:"viewer_count"`
	PublisherId primitive.ObjectID `bson:"publisher_id" json:"publisher_id"`

	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

func NewLiveStream(name string, publisherId primitive.ObjectID) *LiveStream {
	streamId, _ := uuid.NewV6()

	return &LiveStream{
		Name:        name,
		ViewerCount: 0,
		StreamKey:   streamId.String(),
		PublisherId: publisherId,

		CreatedAt: time.Now(),
	}
}
