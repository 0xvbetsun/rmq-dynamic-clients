package rpc

import (
	"bytes"
	"encoding/gob"

	"github.com/streadway/amqp"
	"github.com/vbetsun/rmq-dynamic-clients/configs"
)

type codec struct {
	conn     *amqp.Connection
	ch       *amqp.Channel
	routing  string //routing key
	codec    EncodingCodec
	delivery amqp.Delivery
	msg      <-chan amqp.Delivery
}

type CodecDeps struct {
	Config *configs.Config
	Conn   *amqp.Connection
	Ch     *amqp.Channel
	Queue  amqp.Queue
	Codec  EncodingCodec
	Msg    <-chan amqp.Delivery
}

// EncodingCodec implements marshaling and unmarshaling of seralized data.
type EncodingCodec interface {
	Marshal(interface{}) ([]byte, error)
	Unmarshal([]byte, interface{}) error
}

// GobCodec is an EncodingCodec implementation to send/receive Gob data over AMQP.
type GobCodec struct{}

func (GobCodec) Marshal(v interface{}) ([]byte, error) {
	body := new(bytes.Buffer)
	enc := gob.NewEncoder(body)
	err := enc.Encode(v)
	return body.Bytes(), err
}

func (GobCodec) Unmarshal(data []byte, v interface{}) error {
	body := bytes.NewBuffer(data)
	dec := gob.NewDecoder(body)
	err := dec.Decode(v)
	return err
}
