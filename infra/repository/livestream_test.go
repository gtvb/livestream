package repository

import (
	"testing"

	"github.com/gtvb/livestream/utils"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestCreateLiveStream(t *testing.T) {
	container := setupDatabase()
	defer container.Terminate()

	liveStreamRepo := NewLiveStreamRepository(container.Database, utils.LiveStreamCollectionTest)

	publisherID := primitive.NewObjectID()
	insertedID, err := liveStreamRepo.CreateLiveStream("Test Stream", "streamkey-test", publisherID)

	assert.NoError(t, err)
	assert.NotEqual(t, primitive.NilObjectID, insertedID)
}

func TestDeleteLiveStream(t *testing.T) {
	container := setupDatabase()
	defer container.Terminate()

	liveStreamRepo := NewLiveStreamRepository(container.Database, utils.LiveStreamCollectionTest)

	publisherID := primitive.NewObjectID()

	insertedID, err := liveStreamRepo.CreateLiveStream("Test Stream", "streamkey-test", publisherID)
	assert.NoError(t, err)

	err = liveStreamRepo.DeleteLiveStream(insertedID.(primitive.ObjectID))
	assert.NoError(t, err)
}

func TestDeleteLiveStreamsByPublisher(t *testing.T) {
	container := setupDatabase()
	defer container.Terminate()

	liveStreamRepo := NewLiveStreamRepository(container.Database, utils.LiveStreamCollectionTest)
	publisherID := primitive.NewObjectID()

	_, err := liveStreamRepo.CreateLiveStream("Test Stream 1", "streamkey-test", publisherID)
	assert.NoError(t, err)
	_, err = liveStreamRepo.CreateLiveStream("Test Stream 2", "streamkey-test", publisherID)
	assert.NoError(t, err)

	err = liveStreamRepo.DeleteLiveStreamsByPublisher(publisherID)
	assert.NoError(t, err)
}

func TestUpdateLiveStream(t *testing.T) {
	container := setupDatabase()
	defer container.Terminate()

	liveStreamRepo := NewLiveStreamRepository(container.Database, utils.LiveStreamCollectionTest)
	publisherID := primitive.NewObjectID()

	insertedID, err := liveStreamRepo.CreateLiveStream("Test Stream 1", "streamkey-test", publisherID)
	assert.NoError(t, err)

	err = liveStreamRepo.UpdateLiveStream(insertedID.(primitive.ObjectID), bson.M{"name": "(Updated) Live Stream 1", "live_stream_status": true})
	assert.NoError(t, err)

	ls, _ := liveStreamRepo.GetLiveStreamById(insertedID.(primitive.ObjectID))
	assert.Equal(t, "(Updated) Live Stream 1", ls.Name)
	assert.Equal(t, true, ls.LiveStatus)
}

func TestIncrementLiveStreamUserCount(t *testing.T) {
	container := setupDatabase()
	defer container.Terminate()

	liveStreamRepo := NewLiveStreamRepository(container.Database, utils.LiveStreamCollectionTest)
	publisherID := primitive.NewObjectID()

	insertedID, err := liveStreamRepo.CreateLiveStream("Test Stream", "streamkey-test", publisherID)
	assert.NoError(t, err)

	err = liveStreamRepo.IncrementLiveStreamUserCount(insertedID.(primitive.ObjectID))
	assert.NoError(t, err)
}

func TestDecrementLiveStreamUserCount(t *testing.T) {
	container := setupDatabase()
	defer container.Terminate()

	liveStreamRepo := NewLiveStreamRepository(container.Database, utils.LiveStreamCollectionTest)
	publisherID := primitive.NewObjectID()

	insertedID, err := liveStreamRepo.CreateLiveStream("Test Stream", "streamkey-test", publisherID)
	assert.NoError(t, err)

	err = liveStreamRepo.DecrementLiveStreamUserCount(insertedID.(primitive.ObjectID))
	assert.NoError(t, err)
}

func TestGetLiveStreamById(t *testing.T) {
	container := setupDatabase()
	defer container.Terminate()

	liveStreamRepo := NewLiveStreamRepository(container.Database, utils.LiveStreamCollectionTest)
	publisherID := primitive.NewObjectID()

	insertedID, err := liveStreamRepo.CreateLiveStream("Test Stream", "streamkey-test", publisherID)
	assert.NoError(t, err)

	liveStream, err := liveStreamRepo.GetLiveStreamById(insertedID.(primitive.ObjectID))

	assert.NoError(t, err)
	assert.Equal(t, "Test Stream", liveStream.Name)
}

func TestGetLiveStreamByName(t *testing.T) {
	container := setupDatabase()
	defer container.Terminate()

	liveStreamRepo := NewLiveStreamRepository(container.Database, utils.LiveStreamCollectionTest)

	publisherID := primitive.NewObjectID()
	_, err := liveStreamRepo.CreateLiveStream("Test Stream", "streamkey-test", publisherID)
	assert.NoError(t, err)

	liveStream, err := liveStreamRepo.GetLiveStreamByName("Test Stream")

	assert.NoError(t, err)
	assert.Equal(t, "Test Stream", liveStream.Name)
}

func TestGetAllLiveStreamsByUserId(t *testing.T) {
	container := setupDatabase()
	defer container.Terminate()

	liveStreamRepo := NewLiveStreamRepository(container.Database, utils.LiveStreamCollectionTest)
	publisherID := primitive.NewObjectID()

	_, err := liveStreamRepo.CreateLiveStream("Test Stream 1", "streamkey-test", publisherID)
	assert.NoError(t, err)
	_, err = liveStreamRepo.CreateLiveStream("Test Stream 2", "streamkey-test", publisherID)
	assert.NoError(t, err)

	liveStreams, err := liveStreamRepo.GetAllLiveStreamsByUserId(publisherID)

	assert.NoError(t, err)
	assert.Len(t, liveStreams, 2)
}

func TestGetAllLiveStreams(t *testing.T) {
	container := setupDatabase()
	defer container.Terminate()

	liveStreamRepo := NewLiveStreamRepository(container.Database, utils.LiveStreamCollectionTest)
	publisherID := primitive.NewObjectID()

	_, err := liveStreamRepo.CreateLiveStream("Test Stream 1", "streamkey-test", publisherID)
	assert.NoError(t, err)
	_, err = liveStreamRepo.CreateLiveStream("Test Stream 2", "streamkey-test", publisherID)
	assert.NoError(t, err)

	liveStreams, err := liveStreamRepo.GetAllLiveStreams()

	assert.NoError(t, err)
	assert.Len(t, liveStreams, 2)
}
