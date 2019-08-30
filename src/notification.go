package main

import "go.mongodb.org/mongo-driver/bson/primitive"

type NotificationMessage struct {
	UserId      primitive.ObjectID `json:"userId"`
	WorkspaceId primitive.ObjectID `json:"workspaceId"`
	Amount      uint64             `json:"amount"`
	Timestamp   uint64             `json:"timestamp"`
}
