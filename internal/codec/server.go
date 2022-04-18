package codec

import (
	"errors"
	"net/rpc"
	"sync"
	"sync/atomic"

	"github.com/streadway/amqp"
)

var ErrNoConsumers = errors.New("codec: No consumers in queue")

type route struct {
	messageID string
	routing   string
}

type serverCodec struct {
	*codec
	lock  *sync.RWMutex
	calls map[uint64]route
	seq   uint64
}

func (c *serverCodec) ReadRequestHeader(r *rpc.Request) error {
	c.delivery = <-c.message

	if c.delivery.CorrelationId == "" {
		return errors.New("no routing key in delivery")
	}

	r.Seq = atomic.AddUint64(&c.seq, 1)

	c.lock.Lock()
	c.calls[r.Seq] = route{c.delivery.MessageId, c.delivery.CorrelationId}
	c.lock.Unlock()

	r.ServiceMethod = c.delivery.ReplyTo

	return nil
}

func (c *serverCodec) ReadRequestBody(v interface{}) error {
	if v == nil {
		return nil
	}
	return c.codec.codec.Unmarshal(c.delivery.Body, v)
}

func (c *serverCodec) WriteResponse(resp *rpc.Response, v interface{}) error {
	body, err := c.codec.codec.Marshal(v)
	if err != nil {
		return err
	}

	c.lock.RLock()
	route, ok := c.calls[resp.Seq]
	c.lock.RUnlock()
	if !ok {
		return errors.New("sequence doesn't have a route")
	}

	publishing := amqp.Publishing{
		ContentType:   "application/octet-stream",
		ReplyTo:       resp.ServiceMethod,
		MessageId:     route.messageID,
		CorrelationId: route.routing,
		Body:          body,
	}

	if resp.Error != "" {
		publishing.Headers = amqp.Table{"error": resp.Error}
	}
	return c.channel.Publish(
		"",
		route.routing,
		false,
		false,
		publishing,
	)
}

func (c *serverCodec) Close() error {
	return c.conn.Close()
}

//NewServerCodec returns a new rpc.ClientCodec using AMQP on conn. serverRouting is the routing
//key with with RPC calls are received, encodingCodec is an EncodingCoding implementation.
func NewServerCodec(conn *amqp.Connection, serverRouting string, encodingCodec GobCodec) (rpc.ServerCodec, error) {
	channel, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	queue, err := channel.QueueDeclare(serverRouting, true, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	messages, err := channel.Consume(queue.Name, "", true, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	server := &serverCodec{
		codec: &codec{
			conn:    conn,
			channel: channel,
			routing: queue.Name,

			codec:   encodingCodec,
			message: messages,
		},
		lock:  new(sync.RWMutex),
		calls: make(map[uint64]route),
		seq:   0,
	}

	return server, err
}
