package main

import (
	"github.com/streadway/amqp"
	"log"
)

type QueueMessage struct {
	Payload  []byte `json:"payload"`
	Exchange string `json:"exchange"`
	Route    string `json:"route"`
}

const notificationsExchange = "notifications"
const notificationsRoute = "worker/notifications/checker"

var messagesQueue = make(chan QueueMessage)

func initPublisher() {
	conn := getConnection()

	amqpChannel, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer amqpChannel.Close()

	for msg := range messagesQueue {
		err := amqpChannel.Publish(
			msg.Exchange,
			msg.Route,
			false, // mandatory
			false, // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        msg.Payload,
			})
		if err != nil {
			log.Printf("Publish error: %s", err)
		}
	}
}

func makeNotificationText() {

}
