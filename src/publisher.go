package main

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"log"
)

type QueueMessage struct {
	Type    string              `json:"type"`
	Payload NotificationMessage `json:"payload"`
}

const notificationsExchange = "notify"
const notificationsRoute = "notify/checker"

var messagesQueue = make(chan QueueMessage)

func initPublisher() {
	conn := getConnection()

	amqpChannel, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer amqpChannel.Close()

	for msg := range messagesQueue {
		messageBytes, err := json.Marshal(msg)
		if err != nil {
			return
		}
		err = amqpChannel.Publish(
			notificationsExchange,
			notificationsRoute,
			false, // mandatory
			false, // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        messageBytes,
			})
		if err != nil {
			log.Printf("Publish error: %s", err)
		}
	}
}
