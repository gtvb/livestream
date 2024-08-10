package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/gtvb/livestream/infra/db"
	"github.com/gtvb/livestream/models"
	"go.mongodb.org/mongo-driver/bson"
)

const LiveStreamsCollectionName = "livestreams"

type LiveStreamRepository struct {
	Db *db.Database
}

func NewLiveStreamRepository(db *db.Database) *LiveStreamRepository {
	return &LiveStreamRepository{
		Db: db,
	}
}

func (lr *LiveStreamRepository) CreateLiveStream(name string) (interface{}, error) {
	coll := lr.Db.Collection(LiveStreamsCollectionName)
	doc := models.NewLiveStream(name)

	res, err := coll.InsertOne(context.TODO(), doc)
	if err != nil {
		return nil, err
	}

	return res.InsertedID, nil
}

func (lr *LiveStreamRepository) DeleteLiveStream(id int) (bool, error) {
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

func (lr *LiveStreamRepository) UpdateLiveStreamName(id int, name string) (bool, error) {
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

func (lr *LiveStreamRepository) IncrementLiveStreamUserCount(id int) (bool, error) {
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

func (lr *LiveStreamRepository) DecrementLiveStreamUserCount(id int) (bool, error) {
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
