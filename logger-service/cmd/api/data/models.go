package data

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//create a model with id, name,log,createdat, updatedat and write insert statement

type LogEntry struct {
	Id          string    `bson:"_id,omitempty" json:"id,omitempty"`
	Name        string    `bson:"name" json:"name"`
	Data        string    `bson:"data" json:"data"`
	Created_at  time.Time `bson:"created_at" json:"created_at"`
	Modified_at time.Time `bson:"modified_at" json:"modified_at"`
}

func (l *LogEntry) Insert(entry LogEntry, client *mongo.Client) string {
	collection := client.Database("logs").Collection("logs")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	res, err := collection.InsertOne(ctx, LogEntry{
		Name:        entry.Name,
		Data:        entry.Data,
		Created_at:  time.Now(),
		Modified_at: time.Now(),
	})
	if err != nil {
		fmt.Println("got error while inserting entry", err)
		return ""
	}
	fmt.Println("Inserted result with id", res.InsertedID)

	resString := fmt.Sprintf("Inserted result with id %s", res.InsertedID)
	return resString

}

func (l *LogEntry) GetAll(client *mongo.Client) ([]LogEntry, error) {
	collection := client.Database("logs").Collection("logs")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	opts := options.Find()
	opts.SetSort(bson.D{{"created_at", -1}})
	cur, err := collection.Find(ctx, bson.D{}, opts)
	if err != nil {
		fmt.Println("Got error while fetching records", err)
		return nil, err
	}
	defer cur.Close(ctx)
	var logCollection []LogEntry
	for cur.Next(ctx) {
		var item LogEntry
		err := cur.Decode(&item)
		if err != nil {
			fmt.Println("got error decoding")
			return nil, err
		}
		logCollection = append(logCollection, item)
	}
	return logCollection, nil

}

func (l *LogEntry) GetEntryByID(id string, client *mongo.Client) (*LogEntry, error) {

	collection := client.Database("logs").Collection("logs")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	docid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var item *LogEntry
	err = collection.FindOne(ctx, bson.M{"_id": docid}).Decode(&item)
	if err != nil {
		return nil, err
	}
	return item, nil

}

func (l *LogEntry) DropCollections(client *mongo.Client) error {
	collection := client.Database("logs").Collection("logs")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	if err := collection.Drop(ctx); err != nil {
		return err
	}
	return nil

}

func (l *LogEntry) UpdateCollection(entry LogEntry, client *mongo.Client) error {
	collection := client.Database("logs").Collection("logs")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	primobjId, err := primitive.ObjectIDFromHex(entry.Id)
	if err != nil {
		fmt.Println("Got errr: ", err)
	}
	update := bson.M{"$set": bson.M{"name": entry.Name, "data": entry.Data}}

	if _, err := collection.UpdateOne(ctx, bson.M{"_id": primobjId}, update); err != nil {
		fmt.Println("error in update", err)
		return err
	}

	return nil
}
