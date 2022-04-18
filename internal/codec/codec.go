package codec

import (
	"bytes"
	"encoding/gob"

	"github.com/streadway/amqp"
)

type codec struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	routing string //routing key

	codec    GobCodec
	delivery amqp.Delivery

	message <-chan amqp.Delivery
}

//GobCodec is an EncodingCodec implementation to send/recieve Gob data over AMQP.
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
