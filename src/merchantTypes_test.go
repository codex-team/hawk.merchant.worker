package main

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"testing"
	"time"
)

func TestTransactionFind(t *testing.T) {
	loadEnv()
	database := connectMongo()
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	var tr Transaction
	_ = database.Collection(PaymentTransactionsCollection).FindOne(ctx, bson.D{
		{"orderId", "d25bae71eb2b3870"},
	}).Decode(&tr)
}

func TestUserUpdate(t *testing.T) {
	loadEnv()
	database := connectMongo()
	var tr Transaction
	_, _ = tr.find(database, "d25bae71eb2b3870")
	payment := PaymentAuthorized{
		OrderId:   "64bf62d299bb1990",
		PaymentId: 109847863,
		Status:    "AUTHORIZED",
		Timestamp: 1567166468,
		ErrorCode: 0,
		Amount:    100,
		CardId:    20676544,
		Pan:       "430000******0777",
		ExpDate:   "11/22",
		RebillId:  1567168445546,
	}

	card := UserCard{
		UserId:    tr.UserId,
		CardId:    payment.CardId,
		Pan:       payment.Pan,
		ExpDate:   payment.ExpDate,
		RebillId:  payment.RebillId,
		PaymentId: payment.PaymentId,
	}
	_ = card.insert(database)
}
