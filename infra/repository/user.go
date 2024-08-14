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

const UserCollectionName = "users"

// Repositório de acesso aos dados da entidade `User`.
// Qualquer repositório precisa implementar a interface
// `UserRepositoryInterface` para ser utilizada de forma
// válida pelo servidor HTTP.
type UserRepository struct {
	Db *db.Database
}

func NewUserRepository(db *db.Database) *UserRepository {
	return &UserRepository{
		Db: db,
	}
}

func (ur *UserRepository) CreateUser(name, email, password string) (interface{}, error) {
	coll := ur.Db.Collection(UserCollectionName)
	doc := models.NewUser(name, email, password)

	id, err := coll.InsertOne(context.TODO(), doc)
	if err != nil {
		return primitive.NilObjectID, err
	}

	return id.InsertedID, nil
}

func (ur *UserRepository) DeleteUser(id primitive.ObjectID) (bool, error) {
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

func (ur *UserRepository) UpdateUserName(id primitive.ObjectID, name string) (bool, error) {
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

func (ur *UserRepository) UpdateUserEmail(id primitive.ObjectID, email string) (bool, error) {
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

func (ur *UserRepository) UpdateUserPassword(id primitive.ObjectID, password string) (bool, error) {
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

func (ur *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	coll := ur.Db.Collection(UserCollectionName)
	filter := bson.M{"email": email}

	res := coll.FindOne(context.TODO(), filter)

	var user models.User
	err := res.Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (ur *UserRepository) GetUserById(id primitive.ObjectID) (*models.User, error) {
	coll := ur.Db.Collection(UserCollectionName)
	filter := bson.M{"_id": id}

	res := coll.FindOne(context.TODO(), filter)

	var user models.User
	err := res.Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (ur *UserRepository) GetAllUsers() ([]*models.User, error) {
	coll := ur.Db.Collection(UserCollectionName)
	filter := bson.D{}

	cursor, err := coll.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	var users []*models.User
	if err = cursor.All(context.TODO(), &users); err != nil {
		return nil, err
	}

	return users, nil
}
