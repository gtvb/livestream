package repository

import (
	"context"
	"testing"

	"github.com/gtvb/livestream/utils"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func cleanUpUserCollection(t *testing.T) {
	t.Cleanup(func() {
		// Drop the collection after each test
		err := container.Database.Collection(utils.UserCollectionTest).Drop(context.Background())
		if err != nil {
			t.Fatalf("failed to clean collection %s: %v", utils.UserCollectionTest, err)
		}
	})
}

func TestCreateUser(t *testing.T) {
	container := setupDatabase()
	defer container.Terminate()

	userRepo := NewUserRepository(container.Database, utils.UserCollectionTest)
	insertedID, err := userRepo.CreateUser("johndoe", "johndoe@example.com", "password123")

	assert.NoError(t, err)
	assert.NotEqual(t, primitive.NilObjectID, insertedID)
}

func TestDeleteUser(t *testing.T) {
	container := setupDatabase()
	defer container.Terminate()

	userRepo := NewUserRepository(container.Database, utils.UserCollectionTest)

	insertedID, err := userRepo.CreateUser("johndoe", "johndoe@example.com", "password123")
	assert.NoError(t, err)
	assert.NotEqual(t, primitive.NilObjectID, insertedID)

	err = userRepo.DeleteUser(insertedID.(primitive.ObjectID))
	assert.NoError(t, err)
}

func TestUpdateUser(t *testing.T) {
	container := setupDatabase()
	defer container.Terminate()

	userRepo := NewUserRepository(container.Database, utils.UserCollectionTest)
	insertedID, _ := userRepo.CreateUser("johndoe", "johndoe@example.com", "password123")

	err := userRepo.UpdateUser(insertedID.(primitive.ObjectID), bson.M{"email": "johndoe@new.example.com"})
	assert.Equal(t, nil, err)

	user, _ := userRepo.GetUserById(insertedID.(primitive.ObjectID))
	assert.Equal(t, "johndoe@new.example.com", user.Email)
	assert.Equal(t, "johndoe", user.Username)
}

func TestUpdateAddToFollowList(t *testing.T) {
	container := setupDatabase()
	defer container.Terminate()

	userRepo := NewUserRepository(container.Database, utils.UserCollectionTest)

	insertedIDFollower, err := userRepo.CreateUser("johndoe", "johndoe@example.com", "password123")
	assert.NoError(t, err)

	insertedIDFollowing, err := userRepo.CreateUser("alice", "alice@example.com", "password456")
	assert.NoError(t, err)

	err = userRepo.UpdateUserAddToFollowList(insertedIDFollower.(primitive.ObjectID), insertedIDFollowing.(primitive.ObjectID))
	assert.NoError(t, err)

	userAfterUpdate, err := userRepo.GetUserById(insertedIDFollower.(primitive.ObjectID))
	assert.NoError(t, err)
	assert.Contains(t, userAfterUpdate.Following, insertedIDFollowing)
}

func TestUpdateRemoveFromFollowList(t *testing.T) {
	container := setupDatabase()
	defer container.Terminate()

	userRepo := NewUserRepository(container.Database, utils.UserCollectionTest)

	insertedIDFollower, err := userRepo.CreateUser("johndoe", "johndoe@example.com", "password123")
	assert.NoError(t, err)

	insertedIDFollowing, err := userRepo.CreateUser("alice", "alice@example.com", "password456")
	assert.NoError(t, err)

	userBeforeUpdate, err := userRepo.GetUserById(insertedIDFollower.(primitive.ObjectID))
	assert.NoError(t, err)
	assert.NotContains(t, userBeforeUpdate.Following, insertedIDFollowing)

	err = userRepo.UpdateUserRemoveFromFollowList(insertedIDFollower.(primitive.ObjectID), insertedIDFollowing.(primitive.ObjectID))
	assert.NoError(t, err)

	userAfterUpdate, err := userRepo.GetUserById(insertedIDFollower.(primitive.ObjectID))
	assert.NoError(t, err)
	assert.NotContains(t, userAfterUpdate.Following, insertedIDFollowing)
}

func TestGetUserByUsername(t *testing.T) {
	container := setupDatabase()
	defer container.Terminate()

	userRepo := NewUserRepository(container.Database, utils.UserCollectionTest)

	_, err := userRepo.CreateUser("johndoe", "johndoe@example.com", "password123")
	assert.NoError(t, err)

	user, err := userRepo.GetUserByUsername("johndoe")

	assert.NoError(t, err)
	assert.Equal(t, "johndoe", user.Username)
	assert.Equal(t, "johndoe@example.com", user.Email)
}

func TestGetAllUsers(t *testing.T) {
	container := setupDatabase()
	defer container.Terminate()

	userRepo := NewUserRepository(container.Database, utils.UserCollectionTest)

	_, err := userRepo.CreateUser("johndoe", "johndoe@example.com", "password123")
	assert.NoError(t, err)

	_, err = userRepo.CreateUser("janedoe", "janedoe@example.com", "password123")
	assert.NoError(t, err)

	users, err := userRepo.GetAllUsers()

	assert.NoError(t, err)
	assert.Len(t, users, 2)
}
