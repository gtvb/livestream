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

func (lr *LiveStreamRepository) DeleteLiveStream(id primitive.ObjectID) error {
	coll := lr.Db.Collection(LiveStreamsCollectionName)
	filter := bson.M{"_id": id}

	res, err := coll.DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}

	if res.DeletedCount != 1 {
		return fmt.Errorf("expected one document to be deleted, got %d", res.DeletedCount)
	}

	return nil
}

func (lr *LiveStreamRepository) DeleteLiveStreamsByPublisher(id primitive.ObjectID) error {
	coll := lr.Db.Collection(LiveStreamsCollectionName)
	filter := bson.M{"publisher_id": id}

	_, err := coll.DeleteMany(context.TODO(), filter)
	if err != nil {
		return err
	}

	return nil
}

func (lr *LiveStreamRepository) updateLiveStream(id primitive.ObjectID, updateQuery primitive.M) error {
	coll := lr.Db.Collection(LiveStreamsCollectionName)

	res, err := coll.UpdateByID(context.TODO(), id, updateQuery)
	if err != nil {
		return err
	}

	if res.MatchedCount != 1 {
		return fmt.Errorf("no match for _id %d", id)
	}

	if res.ModifiedCount != 1 {
		return fmt.Errorf("expected one document to be updated, got %d", res.ModifiedCount)
	}

	return nil
}

func (lr *LiveStreamRepository) UpdateLiveStreamSetStatus(id primitive.ObjectID, status bool) error {
	update := bson.M{"$set": bson.M{"live_stream_status": status, "updated_at": time.Now()}}
    return lr.updateLiveStream(id, update)
}

func (lr *LiveStreamRepository) UpdateLiveStreamName(id primitive.ObjectID, name string) error {
	update := bson.M{"$set": bson.M{"name": name, "updated_at": time.Now()}}
    return lr.updateLiveStream(id, update)
}

func (lr *LiveStreamRepository) IncrementLiveStreamUserCount(id primitive.ObjectID) error {
	update := bson.M{
		"$set": bson.M{"updated_at": time.Now()},
		"$inc": bson.M{"viewer_count": 1},
	}
    return lr.updateLiveStream(id, update)
}


func (lr *LiveStreamRepository) DecrementLiveStreamUserCount(id primitive.ObjectID) error {
	update := bson.M{
		"$set": bson.M{"updated_at": time.Now()},
		"$dec": bson.M{"viewer_count": 1},
	}
    return lr.updateLiveStream(id, update)
}

func (lr *LiveStreamRepository) getLiveStreamByParam(fieldName string, param any) (*models.LiveStream, error) {
	var liveStream models.LiveStream
	coll := lr.Db.Collection(LiveStreamsCollectionName)

    filter := bson.M{fieldName: param}

	res := coll.FindOne(context.TODO(), filter)
    err := res.Decode(&liveStream)
	if err != nil {
		return nil, err
	}

	return &liveStream, nil
}

func (lr *LiveStreamRepository) getLiveStreamByParamBatch(filter primitive.M) ([]*models.LiveStream, error) {
	var liveStreams []*models.LiveStream
	coll := lr.Db.Collection(LiveStreamsCollectionName)

	cursor, err := coll.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

    err = cursor.Decode(&liveStreams)
	if err != nil {
		return nil, err
	}

	return liveStreams, nil
}

func (lr *LiveStreamRepository) GetLiveStreamById(id primitive.ObjectID) (*models.LiveStream, error) {
    return lr.getLiveStreamByParam("_id", id)
}

func (lr *LiveStreamRepository) GetLiveStreamByName(name string) (*models.LiveStream, error) {
    return lr.getLiveStreamByParam("name", name)
}

func (lr *LiveStreamRepository) GetLiveStreamByStreamKey(key string) (*models.LiveStream, error) {
    return lr.getLiveStreamByParam("stream_key", key)
}

func (lr *LiveStreamRepository) GetAllLiveStreamsByUserId(id primitive.ObjectID) ([]*models.LiveStream, error) {
    return lr.getLiveStreamByParamBatch(bson.M{"publisher_id": id})
}

// Método genérico, pode ser substituído por uma busca mais específica
func (lr *LiveStreamRepository) GetAllLiveStreams() ([]*models.LiveStream, error) {
    return lr.getLiveStreamByParamBatch(bson.M{})
}
