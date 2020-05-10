package main

import (
	"log"
	"time"

	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func drain(connString string, queue *string, concurrency *int, dwell *int) {
	log.Printf("Draining messages from queue %s", *queue)

	conn, err := amqp.Dial(connString)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		*queue, // name
		false,  // durable
		false,  // delete when unused
		false,  // exclusive
		false,  // no-wait
		nil,    // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.Qos(
		*concurrency, // prefetch count
		0,            // prefetch size
		false,        // global
	)
	failOnError(err, "Failed to set QoS")

	msgs, err := ch.Consume(
		q.Name,                 // queue
		"rabbit-drainer:drain", // consumer
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	failOnError(err, "Failed to register a consumer")

	whenAllIsDone := make(chan bool)

	// Background loop to fetch messages one by one
	go func() {
		for msg := range msgs {
			go func(msg amqp.Delivery) {
				log.Printf("Received message %s", msg.MessageId)
				if *dwell > 0 {
					time.Sleep(time.Duration(*dwell) * time.Millisecond)
				}
				msg.Ack(false)
			}(msg)
		}
	}()

	// Background loop to check if the queue is empty yet
	go func() {
		for {
			time.Sleep(2 * time.Second)
			q, err := ch.QueueInspect(*queue)
			failOnError(err, "Failed to inspect queue")

			log.Printf("Messages remaining: %d", q.Messages)

			if q.Messages == 0 {
				log.Printf("Drained all messages")
				whenAllIsDone <- true
			}
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")

	<-whenAllIsDone
}
