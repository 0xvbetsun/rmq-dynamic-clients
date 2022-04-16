package main

import (
	"bufio"
	"log"
	"os"

	"github.com/streadway/amqp"
)

func main() {
	amqpURL := os.Getenv("AMQP_SERVER_URL")
	amqpQueue := os.Getenv("AMQP_QUEUE_NAME")
	// Create a new RabbitMQ connection.
	connectRabbitMQ, err := amqp.Dial(amqpURL)
	if err != nil {
		panic(err)
	}
	defer connectRabbitMQ.Close()

	// Let's start by opening a channel to our RabbitMQ
	// instance over the connection we have already
	// established.
	channelRabbitMQ, err := connectRabbitMQ.Channel()
	if err != nil {
		panic(err)
	}
	defer channelRabbitMQ.Close()

	// With the instance and declare Queues that we can
	// publish and subscribe to.
	_, err = channelRabbitMQ.QueueDeclare(
		amqpQueue, // queue name
		true,      // durable
		false,     // auto delete
		false,     // exclusive
		false,     // no wait
		nil,       // arguments
	)
	if err != nil {
		panic(err)
	}

	// Make a channel to receive messages into infinite loop.
	forever := make(chan bool)
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		message := amqp.Publishing{
			ContentType: "text/plain",
			Body:        scanner.Bytes(),
		}

		// Attempt to publish a message to the queue.
		if err := channelRabbitMQ.Publish(
			"",              // exchange
			"QueueService1", // queue name
			false,           // mandatory
			false,           // immediate
			message,         // message to publish
		); err != nil {
			log.Println(err)
		}
	}

	<-forever
}
