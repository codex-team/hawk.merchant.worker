package main

type UserBankCard struct {
	UserId     string `json:"userId"`
	CardId     string `json:"cardId"`
	Pan        string `json:"pan"`
	ExpDate    string `json:"expDate"`
	CardType   string `json:"cardType"`
	CardNumber string `json:"cardNumber"`
}

type UserPaymentData struct {
	UserId string         `json:"userId"`
	Email  string         `json:"email"`
	Phone  string         `json:"phone"`
	Cards  []UserBankCard `json:"cards"`
}
