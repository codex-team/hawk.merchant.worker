package main

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"time"
)

const WorkspacesCollection = "workspaces"

func updateWorkspaceBalance(database *mongo.Database, workspaceId string, amount uint32) error {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	workspaceObjectId, err := primitive.ObjectIDFromHex(workspaceId)
	if err != nil {
		log.Printf("[MongoDB] invalid workspaceId: %s", workspaceId)
		return err
	}

	result, err := database.Collection(WorkspacesCollection).UpdateOne(ctx, bson.D{
		{"_id", workspaceObjectId},
	}, bson.D{
		{"$inc", bson.D{
			{
				"balance", amount},
		}},
	})
	if err != nil {
		log.Printf("[MongoDB] updateWorkspaceBalance error: %s", err)
		return err
	}
	if result.ModifiedCount != 1 {
		msg := fmt.Sprintf("[MongoDB] workspace not modified during updateWorkspaceBalance: %s", workspaceId)
		log.Printf(msg)
		return errors.New(msg)
	}
	return nil
}
