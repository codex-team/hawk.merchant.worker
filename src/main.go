package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/codex-team/tinkoff.api.golang"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/streadway/amqp"
)

var (
	tinkoffTerminalKey,
	tinkoffSecretKey,
	amqpURL,
	mongoURL string
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	tinkoffTerminalKey = os.Getenv("TINKOFF_TERMINAL_KEY")
	tinkoffSecretKey = os.Getenv("TINKOFF_SECRET_KEY")
	amqpURL = os.Getenv("AMQP_URL")
	mongoURL = os.Getenv("MONGO_URL")
}

func getConnection() *amqp.Connection {
	conn, err := amqp.Dial(amqpURL)
	failOnError(err, "Failed to connect to RabbitMQ")
	return conn
}

func handleQueue(queueName string, handler func([]byte, *mongo.Database)) {
	conn := getConnection()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.QueueBind(
		q.Name,
		queueName,
		"merchant",
		false,
		nil,
	)
	failOnError(err, "Failed to bind a queue")

	msgs, err := ch.Consume(
		queueName, // queue
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to register a consumer")

	mongoConnection := connectMongo()

	for d := range msgs {
		handler(d.Body, mongoConnection)
	}
}

func initialize() {
	loadEnv()
	conn := getConnection()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"merchant",
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to declare an exchange")
}

func main() {
	forever := make(chan bool)
	initialize()
	go initPublisher()

	go handleQueue("merchant/initialized", func(body []byte, database *mongo.Database) {
		log.Printf("merchant/initialized: %s", body)

		var data PaymentInitialized
		if err := json.Unmarshal(body, &data); err != nil {
			log.Printf("Error: %s", err)
			return
		}

		newPayment := PaymentRequest{
			PaymentId:   data.PaymentId,
			OrderId:     data.OrderId,
			PaymentURL:  data.PaymentURL,
			WorkspaceId: data.WorkspaceId,
			UserId:      data.UserId,
			Status:      data.Status,
			Timestamp:   data.Timestamp,
		}

		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
		res, err := database.Collection("paymentRequests").InsertOne(ctx, newPayment)
		if err != nil {
			log.Printf("MongoDB save error: %s", err)
			return
		}

		log.Printf("payment saved: ID=%s", res.InsertedID)
	})

	go handleQueue("merchant/authorized", func(body []byte, database *mongo.Database) {
		log.Printf("merchant/authorized: %s", body)

		var data PaymentAuthorized
		if err := json.Unmarshal(body, &data); err != nil {
			log.Printf("Error: %s", err)
			return
		}

		newPayment := PaymentRequest{
			PaymentId: data.PaymentId,
			OrderId:   data.OrderId,
			Status:    data.Status,
			ErrorCode: data.ErrorCode,
			Amount:    data.Amount,
			CardId:    data.CardId,
			Pan:       data.Pan,
			ExpDate:   data.ExpDate,
			Timestamp: data.Timestamp,
			RebillId:  data.RebillId,
		}

		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
		res, err := database.Collection("paymentRequests").InsertOne(ctx, newPayment)
		if err != nil {
			log.Printf("MongoDB save error: %s", err)
			return
		}

		log.Printf("payment saved: ID=%s", res.InsertedID)

		_, err = confirm(tinkoff.ConfirmRequest{
			BaseRequest: tinkoff.BaseRequest{
				tinkoffTerminalKey, tinkoffSecretKey,
			},
			PaymentID: newPayment.PaymentId,
			Amount:    newPayment.Amount,
		})
		if err != nil {
			return
		}

		messagesQueue <- QueueMessage{
			Exchange: notificationsExchange,
			Route:    notificationsRoute,
			Payload:  []byte(fmt.Sprintf("Payment confirmed: %d kopecs", data.Amount)),
		}
	})

	go handleQueue("merchant/confirmed", func(body []byte, database *mongo.Database) {
		log.Printf("confirmed: %s", body)
	})

	<-forever
}
