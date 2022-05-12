package amqp

import (
	"errors"
	"log"

	"github.com/streadway/amqp"
	"github.com/vbetsun/rmq-dynamic-clients/configs"
)

type Client struct {
	Conn  *amqp.Connection
	Ch    *amqp.Channel
	Queue amqp.Queue
	Msg   <-chan amqp.Delivery
}

// NewClient returns channel for listening messages
func NewClient(cfg *configs.Config) (*Client, error) {
	conn, err := amqp.Dial(cfg.AMQP.Url)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	queue, err := ch.QueueDeclare(
		"",
		cfg.AMQP.Durable,
		cfg.AMQP.AutoDelete,
		cfg.AMQP.Exclusive,
		cfg.AMQP.NoWait,
		nil,
	)
	if err != nil {
		return nil, err
	}

	sQueue, err := ch.QueueDeclare(
		cfg.AMQP.Queue,
		cfg.AMQP.Durable,
		cfg.AMQP.AutoDelete,
		cfg.AMQP.Exclusive,
		cfg.AMQP.NoWait,
		nil,
	)
	if err != nil {
		return nil, err
	}
	if sQueue.Consumers == 0 {
		return nil, errors.New("no consumers in queue")
	}
	msg, err := ch.Consume(
		queue.Name,
		"",
		cfg.AMQP.AutoAck,
		cfg.AMQP.Exclusive,
		cfg.AMQP.NoLocal,
		cfg.AMQP.NoWait,
		nil,
	)
	if err != nil {
		return nil, err
	}
	return &Client{
		Conn:  conn,
		Ch:    ch,
		Queue: queue,
		Msg:   msg,
	}, nil
}
