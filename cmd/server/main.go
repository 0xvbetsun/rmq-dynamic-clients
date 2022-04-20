// Entry point for application's server
package main

import (
	"log"
	"os"

	"github.com/streadway/amqp"
	"github.com/vbetsun/rmq-dynamic-clients/configs"
	"github.com/vbetsun/rmq-dynamic-clients/internal/rpc"
)

func main() {
	file, err := os.OpenFile("./logs/server.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(file)

	cfg, err := configs.New()
	if err != nil {
		log.Fatal(err)
	}
	conn, err := amqp.Dial(cfg.AMQP.Url)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	serverCodec, err := rpc.NewServerCodec(conn, cfg, rpc.GobCodec{})
	if err != nil {
		log.Fatal(err)
	}
	if err = rpc.Register(rpc.NewItems()); err != nil {
		log.Fatal(err)
	}
	log.Println("Server is running")
	rpc.ServeCodec(serverCodec)
}
