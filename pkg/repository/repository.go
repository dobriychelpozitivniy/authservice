package repository

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"mailer-auth/pkg/models"
	"time"
)

type ConfigRepository struct {
	DBHost     string
	DBUsername string
	DBPassword string
	DBPort     string
	DBTimeout  int
	DBName     string
	DBCollection string
}

type Users interface {
	Read(ctx context.Context, username string) (*models.User, error)
	Create(ctx context.Context, username string, hash string) (interface{}, error)
	CreateTestUser() error
}

type Repository struct {
	Users
}

func NewRepository(cfg *ConfigRepository) (*Repository, error) {
	db, err := connectMongo(cfg)
	if err != nil {
		return nil, err
	}

	return &Repository{
		Users: NewUsersRepository(db, cfg.DBCollection),
	}, nil
}

func connectMongo(c *ConfigRepository) (*mongo.Database, error) {
	connPattern := "mongodb://%v:%v@%v:%v"
	if c.DBUsername == "" {
		connPattern = "mongodb://%s%s%v:%v"
	}

	clientURL := fmt.Sprintf(connPattern,
		c.DBUsername,
		c.DBPassword,
		c.DBHost,
		c.DBPort,
	)
	clientOptions := options.Client().ApplyURI(clientURL)

	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(c.DBTimeout)*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	zap.S().Info("Mongodb has been connected")

	return client.Database(c.DBName), nil
}