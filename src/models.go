package main

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"time"
)

type UserCard struct {
	UserId    string `bson:"userId"`
	CardId    uint32 `bson:"cardId"`
	Pan       string `bson:"pan"`
	ExpDate   string `bson:"expDate"`
	RebillId  uint64 `bson:"rebillId"`
	PaymentId uint64 `bson:"paymentId"`
}

func (uc *UserCard) insert(database *mongo.Database) error {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	var existingCard UserCard
	err := database.Collection(UserCardCollection).FindOne(ctx, bson.D{
		{"userId", uc.UserId},
		{"cardId", uc.CardId},
	}).Decode(&existingCard)
	if err == mongo.ErrNoDocuments {
		res, insertErr := database.Collection(UserCardCollection).InsertOne(ctx, uc)
		if insertErr != nil {
			log.Printf("[MongoDB] New card insert error: %s", err)
			return insertErr
		}
		log.Printf("[MongoDB] Link new card (%d) for user (%s): %s", uc.CardId, uc.UserId, res.InsertedID)
		return nil
	}
	if err != nil {
		log.Printf("[MongoDB] UserCard find error: %s", err)
		return err
	}

	log.Printf("UserCard already exists: %v\n", existingCard)
	log.Printf("New card from bank: %v\n", uc)
	return nil
}
