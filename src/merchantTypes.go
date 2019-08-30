package main

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"time"
)

const PaymentLogsCollection = "paymentLogs"
const PaymentTransactionsCollection = "paymentTransactions"
const TransactionConfirmed = "CONFIRM"
const UsersCollection = "users"
const UserCardCollection = "userCard"

type PaymentInitialized struct {
	PaymentURL    string `json:"paymentURL"`
	TransactionId string `json:"id"`
	UserId        string `json:"userId"`
	WorkspaceId   string `json:"workspaceId"`
	Amount        uint64 `json:"amount"`
	OrderId       string `json:"orderId"`
	PaymentId     uint64 `json:"paymentId,string"`
	Timestamp     uint32 `json:"timestamp"`
}

type PaymentAuthorized struct {
	OrderId   string `json:"orderId"`
	PaymentId uint64 `json:"paymentId"`
	Status    string `json:"status"`
	Timestamp uint32 `json:"timestamp"`
	ErrorCode int    `json:"errorCode,string"`
	Amount    uint32 `json:"amount"`
	CardId    uint32 `json:"cardId"`
	Pan       string `json:"pan"`
	ExpDate   string `json:"expDate"`
	RebillId  uint64 `json:"rebillId"`
}

type Transaction struct {
	TransactionId string `bson:"id"`
	PaymentId     string `bson:"paymentId"`
	OrderId       string `bson:"orderId"`
	UserId        string `bson:"userId"`
	Amount        uint32 `bson:"amount"`
	WorkspaceId   string `bson:"workspaceId"`
	Timestamp     uint32 `bson:"timestamp"`
	Status        string `bson:"status"`
}

type PaymentLog struct {
	OrderId     string    `bson:"orderId"`
	PaymentId   uint64    `bson:"paymentId"`
	PaymentURL  string    `bson:"paymentURL,omitempty"`
	WorkspaceId string    `bson:"workspaceId,omitempty"`
	UserId      string    `bson:"userId,omitempty"`
	Timestamp   time.Time `bson:"timestamp"`
	Status      string    `bson:"status"`
	RebillId    uint64    `bson:"rebillid,omitempty"`
	ErrorCode   int       `bson:"errorCode"`
	Amount      uint32    `bson:"amount,omitempty"`
	CardId      uint32    `bson:"cardId,omitempty"`
	Pan         string    `bson:"pan,omitempty"`
	ExpDate     string    `bson:"expDate,omitempty"`
}

func (pl *PaymentInitialized) save(database *mongo.Database) (*mongo.InsertOneResult, error) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	res, err := database.Collection(PaymentLogsCollection).InsertOne(ctx, pl)
	if err != nil {
		log.Printf("[MongoDB] PaymentInitialized saving error: %s", err)
		return nil, err
	}
	return res, nil
}

func (pl *PaymentAuthorized) save(database *mongo.Database) (*mongo.InsertOneResult, error) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	res, err := database.Collection(PaymentLogsCollection).InsertOne(ctx, pl)
	if err != nil {
		log.Printf("[MongoDB] PaymentAuthorized saving error: %s", err)
		return nil, err
	}
	return res, nil
}

func (tr *Transaction) find(database *mongo.Database, orderId string) (bool, error) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	var t Transaction
	err := database.Collection(PaymentTransactionsCollection).FindOne(ctx, bson.D{
		{"orderId", orderId},
	}).Decode(&t)
	if err == mongo.ErrNoDocuments {
		return false, nil
	}
	if err != nil {
		log.Printf("[MongoDB] Transaction find error: %s", err)
		return false, err
	}
	*tr = t
	return true, nil
}

func (tr *Transaction) update(database *mongo.Database, orderId string) error {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, err := database.Collection(PaymentTransactionsCollection).UpdateOne(ctx, bson.D{
		{"orderId", orderId},
	}, bson.D{
		{"$set", bson.D{
			{
				"status", TransactionConfirmed},
		}},
	})
	if err != nil {
		log.Printf("[MongoDB] Transaction update error: %s", err)
		return err
	}
	if result.ModifiedCount != 1 {
		msg := fmt.Sprintf("[MongoDB] Transaction not found: %s", orderId)
		log.Printf(msg)
		return errors.New(msg)
	}
	return nil
}

func (tr *Transaction) save(database *mongo.Database) error {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	_, err := database.Collection(UserCardCollection).InsertOne(ctx, tr)
	if err != nil {
		return err
	}

	return nil
}
