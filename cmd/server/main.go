// Entry point for application's server
package main

import (
	"log"
	"os"

	"github.com/vbetsun/rmq-dynamic-clients/configs"
	"github.com/vbetsun/rmq-dynamic-clients/internal/service"
	"github.com/vbetsun/rmq-dynamic-clients/internal/storage"
	"github.com/vbetsun/rmq-dynamic-clients/internal/transport/amqp"
	"github.com/vbetsun/rmq-dynamic-clients/internal/transport/rpc"
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
	server, err := amqp.NewServer(cfg)
	if err != nil {
		log.Fatal(err)
	}
	serverCodec := rpc.NewServerCodec(rpc.CodecDeps{
		Config: cfg,
		Conn:   server.Conn,
		Ch:     server.Ch,
		Queue:  server.Queue,
		Codec:  rpc.GobCodec{},
		Msg:    server.Msg,
	})
	if err != nil {
		log.Fatal(err)
	}
	storage := storage.NewItems()
	service := service.NewService(storage)
	handler := rpc.NewHandler(service)
	if err = rpc.Register(handler); err != nil {
		log.Fatal(err)
	}
	log.Println("Server is running")
	rpc.ServeCodec(serverCodec)
}
