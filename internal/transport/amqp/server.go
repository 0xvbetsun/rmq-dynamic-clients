package amqp

import (
	"log"

	"github.com/streadway/amqp"
	"github.com/vbetsun/rmq-dynamic-clients/configs"
)

type Server struct {
	Conn  *amqp.Connection
	Ch    *amqp.Channel
	Queue amqp.Queue
	Msg   <-chan amqp.Delivery
}

// NewServer returns a new instance of server
func NewServer(cfg *configs.Config) (*Server, error) {
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

	msg, err := ch.Consume(
		cfg.AMQP.Queue,
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

	return &Server{
		Conn:  conn,
		Ch:    ch,
		Queue: queue,
		Msg:   msg,
	}, err
}
