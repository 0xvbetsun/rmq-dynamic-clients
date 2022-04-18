package main

import (
	"log"
	"net/rpc"
	"os"

	"github.com/streadway/amqp"
	"github.com/vbetsun/rmq-dynamic-clients/internal/codec"
	"github.com/vbetsun/rmq-dynamic-clients/internal/ordered"
)

type Items struct {
	om *ordered.Map
}

func NewItems() *Items {
	return &Items{
		om: ordered.NewMap(),
	}
}

func (i *Items) AddItem(val string, item *string) error {
	*item = i.om.Set(val)
	log.Printf("add item: %v", i.om)
	return nil
}

func (i *Items) GetItem(val string, item *string) error {
	*item = i.om.Get(val)
	log.Printf("get item: %v", i.om)
	return nil
}

func (i Items) GetAllItems(val string, items *[]string) error {
	*items = i.om.Keys()
	log.Printf("get all items: %v", i.om)
	return nil
}

func (i Items) RemoveItem(val string, deleted *bool) error {
	*deleted = i.om.Delete(val)
	log.Printf("remove item: %v", i.om)
	return nil
}

func main() {
	amqpURL := os.Getenv("AMQP_SERVER_URL")
	if amqpURL == "" {
		log.Fatal("AMQP_SERVER_URL was not provided")
	}
	amqpQueue := os.Getenv("AMQP_QUEUE_NAME")
	if amqpQueue == "" {
		log.Fatal("AMQP_QUEUE_NAME was not provided")
	}
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	serverCodec, err := codec.NewServerCodec(conn, amqpQueue, codec.GobCodec{})
	if err != nil {
		log.Fatal(err)
	}
	rpc.Register(NewItems())
	log.Println("Server is running")
	rpc.ServeCodec(serverCodec)
}
