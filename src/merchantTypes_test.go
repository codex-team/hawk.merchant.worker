package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"testing"
	"time"
)

func TestTransactionFind(t *testing.T) {
	loadEnv()
	database := connectMongo()
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	var tr Transaction
	err := database.Collection(PaymentTransactionsCollection).FindOne(ctx, bson.D{
		{"orderId", "adf40f0339775164"},
	}).Decode(&tr)
	fmt.Printf("%v", tr)
	fmt.Printf("%v", err)
}
