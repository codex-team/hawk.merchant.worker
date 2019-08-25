package main

import "time"

type PaymentInitialized struct {
	OrderId     string    `json:"orderId"`
	PaymentId   uint64    `json:"paymentId,string"`
	PaymentURL  string    `json:"paymentURL"`
	WorkspaceId string    `json:"workspaceId"`
	UserId      string    `json:"userId"`
	Status      string    `json:"status"`
	Timestamp   time.Time `json:"timestamp"`
}

type PaymentAuthorized struct {
	OrderId   string    `json:"orderId"`
	PaymentId uint64    `json:"paymentId"`
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	ErrorCode int       `json:"errorCode,string"`
	Amount    uint64    `json:"amount"`
	CardId    int       `json:"cardId"`
	Pan       string    `json:"pan"`
	ExpDate   string    `json:"expDate"`
}
