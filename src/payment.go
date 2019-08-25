package main

import "time"

type PaymentRequest struct {
	OrderId     string    `bson:"orderId"`
	PaymentId   uint64    `bson:"paymentId"`
	PaymentURL  string    `bson:"paymentURL,omitempty"`
	WorkspaceId string    `bson:"workspaceId,omitempty"`
	UserId      string    `bson:"userId,omitempty"`
	Timestamp   time.Time `bson:"timestamp"`
	Status      string    `bson:"status"`
	RebillId    string    `bson:"rebillid,omitempty"`
	ErrorCode   int       `bson:"errorCode"`
	Amount      uint64    `bson:"amount,omitempty"`
	CardId      int       `bson:"cardId,omitempty"`
	Pan         string    `bson:"pan,omitempty"`
	ExpDate     string    `bson:"expDate,omitempty"`
}
