// Package rpc implements all available requests and handler for communication
package rpc

import (
	"errors"
	"net/rpc"
	"strconv"

	"github.com/streadway/amqp" // TODO: remove this dependency
)

type ClientCodec struct {
	*codec
	rKey string // server routing key
}

// NewClientCodec create instance of Client Code
func NewClientCodec(deps CodecDeps) *ClientCodec {
	return &ClientCodec{
		codec: &codec{
			conn:    deps.Conn,
			ch:      deps.Ch,
			routing: deps.Queue.Name,
			codec:   deps.Codec,
			msg:     deps.Msg,
		},
		rKey: deps.Config.AMQP.Queue,
	}
}

// WriteRequest prepears message and send it to the channel
func (c *ClientCodec) WriteRequest(r *rpc.Request, v interface{}) error {
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
func (c *ClientCodec) ReadResponseHeader(r *rpc.Response) error {
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
func (c *ClientCodec) ReadResponseBody(v interface{}) error {
	if v == nil {
		return nil
	}
	return c.codec.codec.Unmarshal(c.delivery.Body, v)
}

// Close requests and waits for the response to close the AMQP connection.
func (c *ClientCodec) Close() error {
	return c.conn.Close()
}

// NewClientWithCodec returns rpc client
func NewClientWithCodec(codec rpc.ClientCodec) *rpc.Client {
	return rpc.NewClientWithCodec(codec)
}

// AddItem sends request over rpc for adding item to the store
func AddItem(c *rpc.Client, arg string) error {
	return c.Call("Items.AddItem", arg, nil)
}

// GetItem sends request over rpc for retrieving item from store
func GetItem(c *rpc.Client, arg string) error {
	return c.Call("Items.GetItem", arg, nil)
}

// GetAllItems sends request over rpc for retrieving all items in the order they were added
func GetAllItems(c *rpc.Client, arg string) error {
	return c.Call("Items.GetAllItems", arg, nil)
}

// RemoveItem sends request over rpc for deleting requested item from store
func RemoveItem(c *rpc.Client, arg string) error {
	return c.Call("Items.RemoveItem", arg, nil)
}

// BuildRouter creates map for instant finding handler to any command
func BuildRouter() map[string]func(*rpc.Client, string) error {
	return map[string]func(*rpc.Client, string) error{
		"AddItem":     AddItem,
		"RemoveItem":  RemoveItem,
		"GetItem":     GetItem,
		"GetAllItems": GetAllItems,
	}
}
