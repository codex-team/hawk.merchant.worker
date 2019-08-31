package main

import (
	"encoding/json"
	"github.com/codex-team/tinkoff.api.golang"
	"github.com/joho/godotenv"
	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"os"
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

	// handle initialized payments
	go handleQueue("merchant/initialized", func(body []byte, database *mongo.Database) {
		log.Printf("merchant/initialized: %s", body)

		var payment Transaction
		if err := json.Unmarshal(body, &payment); err != nil {
			log.Printf("[PaymentInitialized] Unmarshal payload error: %s", err)
			return
		}

		err := payment.save(database)
		if err != nil {
			return
		}
		log.Printf("[initialized] payment saved: orderId=%s", payment.OrderId)
	})

	go handleQueue("merchant/authorized", func(body []byte, database *mongo.Database) {
		log.Printf("merchant/authorized: %s", body)

		var payment PaymentAuthorized
		if err := json.Unmarshal(body, &payment); err != nil {
			log.Printf("[PaymentAuthorized] Unmarshal payload error: %s", err)
			return
		}

		res, err := payment.save(database)
		if err != nil {
			return
		}
		log.Printf("[authorized] payment saved: ID=%s", res.InsertedID)

		// check transaction existence
		var transaction Transaction
		ok, err := transaction.find(database, payment.OrderId)
		if err != nil {
			return
		}
		if !ok {
			log.Printf("[authorized] transaction not found: %s", payment.OrderId)
			return
		}

		log.Printf("found transaction: %v", transaction)

		_, err = confirm(tinkoff.ConfirmRequest{
			BaseRequest: tinkoff.BaseRequest{
				TerminalKey: tinkoffTerminalKey,
				Token:       tinkoffSecretKey,
			},
			PaymentID: uint64(payment.PaymentId),
			Amount:    uint64(payment.Amount),
		})
		if err != nil {
			return
		}

		status := transaction.Status

		if err := transaction.update(database, payment.OrderId); err != nil {
			return
		}

		if status != TransactionSinglePayment {
			card := UserCard{
				UserId:    transaction.UserId,
				CardId:    payment.CardId,
				Pan:       payment.Pan,
				ExpDate:   payment.ExpDate,
				RebillId:  payment.RebillId,
				PaymentId: payment.PaymentId,
			}

			err = card.insert(database)
			failOnError(err, "UserCard saving error during confirmation stage")
		} else {
			err = updateWorkspaceBalance(database, transaction.WorkspaceId, transaction.Amount)
			if err != nil {
				log.Printf("Balance update error: %s", err)
				return
			}
		}

		messagesQueue <- QueueMessage{
			Type: "merchant",
			Payload: NotificationMessage{
				Amount:      transaction.Amount,
				UserId:      transaction.UserId,
				WorkspaceId: transaction.WorkspaceId,
				Timestamp:   transaction.Timestamp,
			},
		}
	})

	go handleQueue("merchant/confirmed", func(body []byte, database *mongo.Database) {
		log.Printf("confirmed: %s", body)
	})

	go handleQueue("merchant/charged", func(body []byte, database *mongo.Database) {
		log.Printf("charged: %s", body)

		var transaction Transaction
		if err := json.Unmarshal(body, &transaction); err != nil {
			log.Printf("Charge unmarshal error: %s", err)
			return
		}

		if err := transaction.save(database); err != nil {
			log.Printf("Transaction save error: %s", err)
			return
		}

		err := updateWorkspaceBalance(database, transaction.WorkspaceId, transaction.Amount)
		if err != nil {
			log.Printf("Balance update error: %s", err)
			return
		}
	})

	log.Printf("Server started:\n\t- AMQP: %s\n\t- MongoDB: %s\n", amqpURL, mongoURL)

	<-forever
}
