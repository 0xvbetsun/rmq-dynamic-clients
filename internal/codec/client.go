package codec

import (
	"errors"
	"net/rpc"
	"strconv"

	"github.com/streadway/amqp"
)

type clientCodec struct {
	*codec
	serverRouting string //server routing key
}

func (c *clientCodec) WriteRequest(r *rpc.Request, v interface{}) error {
	body, err := c.codec.codec.Marshal(v)
	if err != nil {
		return err
	}
	publishing := amqp.Publishing{
		ContentType:   "application/octet-stream",
		ReplyTo:       r.ServiceMethod,
		CorrelationId: c.routing,
		MessageId:     strconv.FormatUint(r.Seq, 10),
		Body:          body,
	}
	return c.channel.Publish("", c.serverRouting, false, false, publishing)
}

func (c *clientCodec) ReadResponseHeader(r *rpc.Response) error {
	c.delivery = <-c.message
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
	r.Seq, err = strconv.ParseUint(c.delivery.MessageId, 10, 64)
	if err != nil {
		return err
	}

	return nil
}

func (c *clientCodec) ReadResponseBody(v interface{}) error {
	if v == nil {
		return nil
	}
	return c.codec.codec.Unmarshal(c.delivery.Body, v)
}

func (c *clientCodec) Close() error {
	return c.conn.Close()
}

//NewClientCodec returns a new rpc.ClientCodec using AMQP on conn. serverRouting is the routing
//key with with RPC calls are sent, it should be the same routing key used with NewServerCodec.
func NewClientCodec(conn *amqp.Connection, serverRouting string, encodingCodec GobCodec) (rpc.ClientCodec, error) {
	channel, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	serverQueue, err := channel.QueueDeclare(serverRouting, true, false, false, false, nil)
	if err != nil {
		return nil, err
	}
	if serverQueue.Consumers == 0 {
		return nil, ErrNoConsumers
	}

	message, err := channel.Consume(serverQueue.Name, "", true, false, false, false, nil)
	client := &clientCodec{
		codec: &codec{
			conn:    conn,
			channel: channel,
			routing: serverRouting,
			codec:   encodingCodec,
			message: message,
		},

		serverRouting: serverRouting,
	}

	return client, err
}
