package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"mailer-auth/pkg/models"
)

type UsersRepository struct {
	db           *mongo.Database
	dbCollection string
}

func NewUsersRepository(db *mongo.Database, dbCollection string) *UsersRepository {
	return &UsersRepository{db: db, dbCollection: dbCollection}
}

// Gets a user(models.User) from the database by his username
func (r *UsersRepository) Read(ctx context.Context, username string) (*models.User, error) {
	filter := bson.D{
		{Key: "username", Value: bson.M{"$eq": username}},
	}

	res := r.db.Collection(r.dbCollection).FindOne(ctx, filter)

	var user *models.User
	err := res.Decode(&user)

	if err != nil {
		return nil, err
	}

	return user, nil
}

// Create user in database
func (r *UsersRepository) Create(ctx context.Context, username, hash string) (interface{}, error) {
	document := &models.User{
		ID:           primitive.ObjectID{},
		Username:     username,
		PasswordHash: hash,
	}

	res, err := r.db.Collection(r.dbCollection).InsertOne(ctx, document)
	if err != nil {
		return nil, err
	}

	return res.InsertedID, nil
}

// Create user in db for testing
func (r *UsersRepository) CreateTestUser() error {
	filter := bson.D{
		{
			Key:   "username",
			Value: bson.M{"$eq": "testuser"},
		},
	}

	res := r.db.Collection(r.dbCollection).FindOne(context.TODO(), filter)

	var user *models.User

	err := res.Decode(&user)
	if err != nil {
		_, err = r.Create(context.TODO(), "testuser", "$2a$12$tc548.1ls5q7Enkgsj5ivuP0LRU1ATyp0TaAWaSWFveZZd59TmeZm")
		if err != nil {
			zap.S().Panicf("Error create testuser in db: %s", err)

			return err
		}
	}

	return nil
}
