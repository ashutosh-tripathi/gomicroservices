package router

import (
	"context"
	"logger-service/cmd/api/data"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type RPCServer struct{}

type RPCPayload struct {
	Name string
	Data string
}

var client *mongo.Client

func (r *RPCServer) LogInfo(payload RPCPayload, resp *string) error {

	collection := client.Database("logs").Collection("logs")
	collection.InsertOne(context.TODO(), data.LogEntry{
		Name:        payload.Name,
		Data:        payload.Data,
		Created_at:  time.Now(),
		Modified_at: time.Now(),
	})
	*resp = "processed payload via RPC" + payload.Name
	return nil
}
