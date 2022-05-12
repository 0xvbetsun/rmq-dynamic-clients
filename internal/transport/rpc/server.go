package rpc

import (
	"errors"
	"net/rpc"
	"sync"
	"sync/atomic"

	"github.com/streadway/amqp" // TODO: remove this dependency
)

type route struct {
	msgID   string
	routing string
}

type serverCodec struct {
	*codec
	lock  *sync.RWMutex
	calls map[uint64]route
	seq   uint64
}

// NewServerCodec returns a new instance of Server Codec
func NewServerCodec(deps CodecDeps) rpc.ServerCodec {
	return &serverCodec{
		codec: &codec{
			conn:    deps.Conn,
			ch:      deps.Ch,
			routing: deps.Queue.Name,
			codec:   deps.Codec,
			msg:     deps.Msg,
		},
		lock:  new(sync.RWMutex),
		calls: make(map[uint64]route),
		seq:   0,
	}
}

// ReadRequestHeader validates headers from request
func (c *serverCodec) ReadRequestHeader(r *rpc.Request) error {
	c.delivery = <-c.msg
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

// ReadRequestBody unmarshals delivered data in body
func (c *serverCodec) ReadRequestBody(v interface{}) error {
	if v == nil {
		return nil
	}
	return c.codec.codec.Unmarshal(c.delivery.Body, v)
}

// WriteResponse prepears message and send it to the channel
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

	msg := amqp.Publishing{
		ContentType:   "application/octet-stream",
		ReplyTo:       resp.ServiceMethod,
		MessageId:     route.msgID,
		CorrelationId: route.routing,
		Body:          body,
	}

	if resp.Error != "" {
		msg.Headers = amqp.Table{"error": resp.Error}
	}
	return c.ch.Publish(
		"",
		route.routing,
		false,
		false,
		msg,
	)
}

// Close requests and waits for the response to close the AMQP connection.
func (c *serverCodec) Close() error {
	return c.conn.Close()
}

// Register publishes the receiver's methods in the rpc.DefaultServer.
func Register(receiver interface{}) error {
	return rpc.Register(receiver)
}

// ServeCodec is like ServeConn but uses the specified codec to
// decode requests and encode responses.
func ServeCodec(codec rpc.ServerCodec) {
	rpc.ServeCodec(codec)
}
