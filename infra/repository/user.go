package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/gtvb/livestream/infra/db"
	"github.com/gtvb/livestream/models"
	"go.mongodb.org/mongo-driver/bson"
)

const UserCollectionName = "users"

type UserRepository struct {
	Db *db.Database
}

func NewUserRepository(db *db.Database) *LiveStreamRepository {
	return &LiveStreamRepository{
		Db: db,
	}
}

func (ur *UserRepository) CreateUser(name, email, password string) (interface{}, error) {
	coll := ur.Db.Collection(UserCollectionName)
	doc := models.NewUser(name, email, password)

	res, err := coll.InsertOne(context.TODO(), doc)
	if err != nil {
		return nil, err
	}

	return res.InsertedID, nil
}

func (ur *UserRepository) DeleteUser(id int) (bool, error) {
	coll := ur.Db.Collection(UserCollectionName)
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

func (ur *UserRepository) UpdateUserName(id int, name string) (bool, error) {
	coll := ur.Db.Collection(UserCollectionName)
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

func (ur *UserRepository) UpdateUserEmail(id int, email string) (bool, error) {
	coll := ur.Db.Collection(UserCollectionName)
	update := bson.M{"$set": bson.M{"email": email, "updated_at": time.Now()}}

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

func (ur *UserRepository) UpdateUserPassword(id int, password string) (bool, error) {
	coll := ur.Db.Collection(UserCollectionName)
	update := bson.M{"$set": bson.M{"password": password, "updated_at": time.Now()}}

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
