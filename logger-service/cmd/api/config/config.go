package config

import (
	"context"
	"errors"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DBclient1 *mongo.Client

type ConfigParam struct {
	DBClient   *mongo.Client
	DBURL      string
	DBUsername string
	DBPassword string
	WebPort    int
	RPCPort    int
	GRPCPort   int
}

func (c *ConfigParam) CreateConnection(ctx context.Context) (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(c.DBURL).SetMaxPoolSize(10)
	clientOptions.SetAuth(options.Credential{
		Username: c.DBUsername,
		Password: c.DBPassword,
	})

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Println("Error encountered while connecting to MongoDB", err)
		return nil, err
	}

	c.DBClient = client
	DBclient1 = client
	return client, nil
}
func (c *ConfigParam) GetConnection() (*mongo.Client, error) {
	if c.DBClient == nil {
		return nil, errors.New("DB client is not initialized")
	} else {
		return c.DBClient, nil
	}
}
func GetClient() (*mongo.Client, error) {
	if DBclient1 == nil {
		return nil, errors.New("DB client is not initialized")
	} else {
		return DBclient1, nil
	}
}
