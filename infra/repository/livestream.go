package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/gtvb/livestream/infra/db"
	"github.com/gtvb/livestream/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const LiveStreamsCollectionName = "livestreams"

// Repositório de acesso aos dados da entidade `LiveStream`.
// Qualquer repositório precisa implementar a interface
// `LiveStreamRepositoryInterface` para ser utilizada de forma
// válida pelo servidor HTTP.
type LiveStreamRepository struct {
	Db *db.Database
}

func NewLiveStreamRepository(db *db.Database) *LiveStreamRepository {
	return &LiveStreamRepository{
		Db: db,
	}
}

func (lr *LiveStreamRepository) CreateLiveStream(name string, publisherId primitive.ObjectID) (interface{}, error) {
	coll := lr.Db.Collection(LiveStreamsCollectionName)
	doc := models.NewLiveStream(name, publisherId)

	res, err := coll.InsertOne(context.TODO(), doc)
	if err != nil {
		return nil, err
	}

	return res.InsertedID, nil
}

func (lr *LiveStreamRepository) DeleteLiveStream(id primitive.ObjectID) (bool, error) {
	coll := lr.Db.Collection(LiveStreamsCollectionName)
	filter := bson.M{"_id": id}

	res, err := coll.DeleteOne(context.TODO(), filter)
	if err != nil {
		return false, err
	}

	if res.DeletedCount != 1 {
		return false, fmt.Errorf("expected one document to be deleted, got %d", res.DeletedCount)
	}

	return true, nil
}

func (lr *LiveStreamRepository) DeleteLiveStreamsByPublisher(id primitive.ObjectID) (bool, error) {
	coll := lr.Db.Collection(LiveStreamsCollectionName)
	filter := bson.M{"publisher_id": id}

	_, err := coll.DeleteMany(context.TODO(), filter)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (lr *LiveStreamRepository) UpdateLiveStreamName(id primitive.ObjectID, name string) (bool, error) {
	coll := lr.Db.Collection(LiveStreamsCollectionName)
	update := bson.M{"$set": bson.M{"name": name, "updated_at": time.Now()}}

	res, err := coll.UpdateByID(context.TODO(), id, update)
	if err != nil {
		return false, err
	}

	if res.MatchedCount != 1 {
		return false, fmt.Errorf("no match for _id %d", id)
	}

	if res.ModifiedCount != 1 {
		return false, fmt.Errorf("expected one document to be updated, got %d", res.ModifiedCount)
	}

	return true, nil
}

func (lr *LiveStreamRepository) IncrementLiveStreamUserCount(id primitive.ObjectID) (bool, error) {
	coll := lr.Db.Collection(LiveStreamsCollectionName)
	update := bson.M{
		"$set": bson.M{"updated_at": time.Now()},
		"$inc": bson.M{"viewer_count": 1},
	}

	res, err := coll.UpdateByID(context.TODO(), id, update)
	if err != nil {
		return false, err
	}

	if res.MatchedCount != 1 {
		return false, fmt.Errorf("no match for _id %d", id)
	}

	if res.ModifiedCount != 1 {
		return false, fmt.Errorf("expected one document to be updated, got %d", res.ModifiedCount)
	}

	return true, nil
}

func (lr *LiveStreamRepository) DecrementLiveStreamUserCount(id primitive.ObjectID) (bool, error) {
	coll := lr.Db.Collection(LiveStreamsCollectionName)
	update := bson.M{
		"$set": bson.M{"updated_at": time.Now()},
		"$dec": bson.M{"viewer_count": 1},
	}

	res, err := coll.UpdateByID(context.TODO(), id, update)
	if err != nil {
		return false, err
	}

	if res.MatchedCount != 1 {
		return false, fmt.Errorf("no match for _id %d", id)
	}

	if res.ModifiedCount != 1 {
		return false, fmt.Errorf("expected one document to be updated, got %d", res.ModifiedCount)
	}

	return true, nil
}

func (lr *LiveStreamRepository) GetLiveStreamByName(name string) (*models.LiveStream, error) {
	coll := lr.Db.Collection(LiveStreamsCollectionName)
	filter := bson.M{"name": name}

	cursor, err := coll.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	var liveStream models.LiveStream
	err = cursor.Decode(&liveStream)
	if err != nil {
		return nil, err
	}

	return &liveStream, nil
}

func (lr *LiveStreamRepository) GetAllLiveStreamsByUserId(id primitive.ObjectID) ([]*models.LiveStream, error) {
	coll := lr.Db.Collection(LiveStreamsCollectionName)
	filter := bson.M{"publisher_id": id}

	cursor, err := coll.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	var liveStreams []*models.LiveStream
	err = cursor.All(context.TODO(), &liveStreams)
	if err != nil {
		return nil, err
	}

	return liveStreams, nil
}

// This is a generic method, just so we can display a
// significant number of streams on the client for
// testing purposes. Later on, we could add tagging
// capabilities, or even expand to more complex searching
// techniques
func (lr *LiveStreamRepository) GetAllLiveStreams() ([]*models.LiveStream, error) {
	coll := lr.Db.Collection(LiveStreamsCollectionName)

	cursor, err := coll.Find(context.TODO(), bson.D{{}})
	if err != nil {
		return nil, err
	}

	var liveStreams []*models.LiveStream
	err = cursor.All(context.TODO(), &liveStreams)
	if err != nil {
		return nil, err
	}

	return liveStreams, nil
}
