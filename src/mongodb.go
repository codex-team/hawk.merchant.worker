package main

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type MongoConnection struct {
	Client   *mongo.Client
	Database *mongo.Database
}

func connectMongo() *mongo.Database {
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURL))
	failOnError(err, "MongoDB client initialization")
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	err = client.Connect(ctx)
	failOnError(err, "MondbDB client connect")
	database := client.Database("hawk")
	return database
}
