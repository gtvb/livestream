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

// Repositório de acesso aos dados da entidade `User`.
// Qualquer repositório precisa implementar a interface
// `UserRepositoryInterface` para ser utilizada de forma
// válida pelo servidor HTTP.
type UserRepository struct {
	userCollectionName string
	Db                 *db.Database
}

func NewUserRepository(db *db.Database, userCollectionName string) *UserRepository {
	return &UserRepository{
		userCollectionName: userCollectionName,
		Db:                 db,
	}
}

func (ur *UserRepository) CreateUser(name, username, email, password string) (interface{}, error) {
	coll := ur.Db.Collection(ur.userCollectionName)
	doc := models.NewUser(name, username, email, password)

	// TODO: verify email against a valid pattern and return error if it exists, implement it

	id, err := coll.InsertOne(context.TODO(), doc)
	if err != nil {
		return primitive.NilObjectID, err
	}

	return id.InsertedID, nil
}

func (ur *UserRepository) DeleteUser(id primitive.ObjectID) error {
	coll := ur.Db.Collection(ur.userCollectionName)
	filter := bson.M{"_id": id}

	_, err := coll.DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}

	return nil
}

func (ur *UserRepository) updateUser(id primitive.ObjectID, updateQuery primitive.M) error {
	coll := ur.Db.Collection(ur.userCollectionName)

	res, err := coll.UpdateByID(context.TODO(), id, updateQuery)
	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		return fmt.Errorf("no match for _id %d", id)
	}

	if res.ModifiedCount != 1 {
		return fmt.Errorf("expected one document to be updated, got %d", res.ModifiedCount)
	}

	return nil
}

func (ur *UserRepository) UpdateUserName(id primitive.ObjectID, name string) error {
	update := bson.M{"$set": bson.M{"name": name, "updated_at": time.Now()}}
	return ur.updateUser(id, update)
}

func (ur *UserRepository) UpdateUserEmail(id primitive.ObjectID, email string) error {
	update := bson.M{"$set": bson.M{"email": email, "updated_at": time.Now()}}
	return ur.updateUser(id, update)
}

func (ur *UserRepository) UpdateUserPassword(id primitive.ObjectID, password string) error {
	update := bson.M{"$set": bson.M{"password": password, "updated_at": time.Now()}}
	return ur.updateUser(id, update)
}

func (ur *UserRepository) getUserByParam(filter primitive.M) (*models.User, error) {
	var user models.User
	coll := ur.Db.Collection(ur.userCollectionName)

	res := coll.FindOne(context.TODO(), filter)
	err := res.Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (ur *UserRepository) GetUserByUsername(username string) (*models.User, error) {
	return ur.getUserByParam(bson.M{"username": username})
}

func (ur *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	return ur.getUserByParam(bson.M{"email": email})
}

func (ur *UserRepository) GetUserById(id primitive.ObjectID) (*models.User, error) {
	return ur.getUserByParam(bson.M{"_id": id})
}

func (ur *UserRepository) GetAllUsers() ([]*models.User, error) {
	coll := ur.Db.Collection(ur.userCollectionName)
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
