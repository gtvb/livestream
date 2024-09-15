package repository

import (
	"context"
	"testing"

	"github.com/gtvb/livestream/utils"
	"github.com/stretchr/testify/assert"
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
	insertedID, err := userRepo.CreateUser("John Doe", "johndoe", "johndoe@example.com", "password123")

	assert.NoError(t, err)
	assert.NotEqual(t, primitive.NilObjectID, insertedID)
}

func TestDeleteUser(t *testing.T) {
	container := setupDatabase()
	defer container.Terminate()

	userRepo := NewUserRepository(container.Database, utils.UserCollectionTest)

	insertedID, err := userRepo.CreateUser("John Doe", "johndoe", "johndoe@example.com", "password123")
	assert.NoError(t, err)
	assert.NotEqual(t, primitive.NilObjectID, insertedID)

	err = userRepo.DeleteUser(insertedID.(primitive.ObjectID))
	assert.NoError(t, err)
}

func TestUpdateUserName(t *testing.T) {
	container := setupDatabase()
	defer container.Terminate()

	userRepo := NewUserRepository(container.Database, utils.UserCollectionTest)

	insertedID, err := userRepo.CreateUser("John Doe", "johndoe", "johndoe@example.com", "password123")
	assert.NoError(t, err)

	err = userRepo.UpdateUserName(insertedID.(primitive.ObjectID), "New Name")
	assert.NoError(t, err)
}

func TestGetUserByUsername(t *testing.T) {
	container := setupDatabase()
	defer container.Terminate()

	userRepo := NewUserRepository(container.Database, utils.UserCollectionTest)

	_, err := userRepo.CreateUser("John Doe", "johndoe", "johndoe@example.com", "password123")
	assert.NoError(t, err)

	user, err := userRepo.GetUserByUsername("johndoe")

	assert.NoError(t, err)
	assert.Equal(t, "John Doe", user.Name)
	assert.Equal(t, "johndoe", user.Username)
	assert.Equal(t, "johndoe@example.com", user.Email)
}

func TestGetAllUsers(t *testing.T) {
	container := setupDatabase()
	defer container.Terminate()

	userRepo := NewUserRepository(container.Database, utils.UserCollectionTest)

	_, err := userRepo.CreateUser("John Doe", "johndoe", "johndoe@example.com", "password123")
	assert.NoError(t, err)

	_, err = userRepo.CreateUser("Jane Doe", "janedoe", "janedoe@example.com", "password123")
	assert.NoError(t, err)

	users, err := userRepo.GetAllUsers()

	assert.NoError(t, err)
	assert.Len(t, users, 2)
}
