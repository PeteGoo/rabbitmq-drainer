package main

import (
	"log"

	"github.com/streadway/amqp"
)

func receive(connString string, queue *string, count *int, durable *bool, verbose *bool) {
	conn, err := amqp.Dial(connString)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		*queue,   // name
		*durable, // durable
		false,    // delete when unused
		false,    // exclusive
		false,    // no-wait
		nil,      // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	failOnError(err, "Failed to set QoS")

	msgs, err := ch.Consume(
		q.Name,                     // queue
		"rabbitmq-drainer:receive", // consumer
		false,                      // auto-ack
		false,                      // exclusive
		false,                      // no-local
		false,                      // no-wait
		nil,                        // args
	)
	failOnError(err, "Failed to register a consumer")

	for i := 0; i < *count; i++ {
		d := <-msgs
		log.Printf("%s", d.Body)
		d.Ack(false)
	}
}
