package main

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

type UserPaymentData struct {
	UserId string     `json:"userId"`
	Email  string     `json:"email"`
	Phone  string     `json:"phone"`
	Cards  []UserCard `json:"cards"`
}

type User struct {
	Id primitive.ObjectID `bson:"_id"`
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

//func (uc *UserCard) attachCard(card UserCard) error {
//	found := false
//	for _, lookupCard := range user.Cards {
//		if lookupCard.CardId == card.CardId {
//			found = true
//			break
//		}
//	}
//	if !found {
//		user.Cards = append(user.Cards, card)
//	}
//
//	return nil
//}
//
//func (user *User) save(database *mongo.Database) error {
//	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
//	result, err := database.Collection(UsersCollection).UpdateOne(ctx, bson.D{
//		{"_id", user.Id},
//	}, bson.D{
//		{"$set", bson.D{
//			{
//				"cards", user.Cards},
//		}},
//	})
//	if err != nil {
//		log.Printf("[MongoDB] User update error: %s", err)
//		return err
//	}
//	if result.ModifiedCount != 1 {
//		msg := fmt.Sprintf("[MongoDB] User not found: %s", user.Id)
//		log.Printf(msg)
//		return errors.New(msg)
//	}
//	return nil
//}
