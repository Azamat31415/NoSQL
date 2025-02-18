package config

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

var client *mongo.Client

func ConnectMongoDB() *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	c, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	if err = c.Ping(ctx, nil); err != nil {
		log.Fatal(err)
	}

	client = c
	return client
}

func GetCollection(dbName, collectionName string) *mongo.Collection {
	return client.Database(dbName).Collection(collectionName)
}
