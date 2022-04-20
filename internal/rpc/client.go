package rpc

import (
	"errors"
	"net/rpc"
	"strconv"

	"github.com/streadway/amqp"
	"github.com/vbetsun/rmq-dynamic-clients/configs"
)

type clientCodec struct {
	*codec
	rKey string // server routing key
}

// WriteRequest prepears message and send it to the channel
func (c *clientCodec) WriteRequest(r *rpc.Request, v interface{}) error {
	body, err := c.codec.codec.Marshal(v)
	if err != nil {
		return err
	}
	msg := amqp.Publishing{
		ContentType:   "application/octet-stream",
		ReplyTo:       r.ServiceMethod,
		CorrelationId: c.routing,
		MessageId:     strconv.FormatUint(r.Seq, 10),
		Body:          body,
	}
	return c.ch.Publish("", c.rKey, false, false, msg)
}

// ReadResponseHeader reads and validates headers from response
func (c *clientCodec) ReadResponseHeader(r *rpc.Response) error {
	c.delivery = <-c.msg
	var err error
	if err = c.delivery.Headers.Validate(); err != nil {
		return errors.New("error while reading body: " + err.Error())
	}

	if err, ok := c.delivery.Headers["error"]; ok {
		errStr, ok := err.(string)
		if !ok {
			return errors.New("error header not a string")
		}
		r.Error = errStr
	}

	r.ServiceMethod = c.delivery.ReplyTo
	seqID, err := strconv.ParseUint(c.delivery.MessageId, 10, 64)
	if err != nil {
		seqID = 0
	}
	r.Seq = seqID
	return nil
}

// ReadResponseBody unmarshals delivered body
func (c *clientCodec) ReadResponseBody(v interface{}) error {
	if v == nil {
		return nil
	}
	return c.codec.codec.Unmarshal(c.delivery.Body, v)
}

// Close requests and waits for the response to close the AMQP connection.
func (c *clientCodec) Close() error {
	return c.conn.Close()
}

// NewClientCodec returns a new rpc.ClientCodec using AMQP on conn. serverRouting is the routing
// key with with RPC calls are sent, it should be the same routing key used with NewServerCodec.
func NewClientCodec(conn *amqp.Connection, cfg *configs.Config, encodingCodec EncodingCodec) (rpc.ClientCodec, error) {
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
	client := &clientCodec{
		codec: &codec{
			conn:    conn,
			ch:      ch,
			routing: queue.Name,
			codec:   encodingCodec,
			msg:     msg,
		},
		rKey: cfg.AMQP.Queue,
	}

	return client, err
}

// NewClientWithCodec returns rpc client
func NewClientWithCodec(codec rpc.ClientCodec) *rpc.Client {
	return rpc.NewClientWithCodec(codec)
}
