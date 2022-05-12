// Entry point for application's clients
package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/vbetsun/rmq-dynamic-clients/configs"
	"github.com/vbetsun/rmq-dynamic-clients/internal/transport/amqp"
	"github.com/vbetsun/rmq-dynamic-clients/internal/transport/rpc"
)

func main() {
	cfg, err := configs.New()
	if err != nil {
		log.Fatal(err)
	}
	amqpClient, err := amqp.NewClient(cfg)
	if err != nil {
		log.Fatal(err)
	}
	cmdRouter := rpc.BuildRouter()
	clientCodec := rpc.NewClientCodec(rpc.CodecDeps{
		Config: cfg,
		Conn:   amqpClient.Conn,
		Ch:     amqpClient.Ch,
		Queue:  amqpClient.Queue,
		Codec:  rpc.GobCodec{},
		Msg:    amqpClient.Msg,
	})
	rpcClient := rpc.NewClientWithCodec(clientCodec)
	scanner := bufio.NewScanner(os.Stdin)
	log.Println("Client is running")
	for scanner.Scan() {
		cmd, arg, err := parseCmd(scanner.Text())
		if err != nil {
			log.Println(err)
			continue
		}
		if handler, ok := cmdRouter[cmd]; ok {
			err = handler(rpcClient, arg)
			if err != nil {
				log.Print(err)
			}
		} else {
			log.Printf("unknown command %s", cmd)
		}
	}
}

func parseCmd(cmd string) (string, string, error) {
	if cmd == "" {
		return "", "", errors.New("command can not be empty string")
	}
	el := strings.Fields(cmd)
	if len(el) > 2 {
		return "", "", fmt.Errorf("illegal number of arguments %d", len(el))
	}
	// example: GetAllItems passes without args
	if len(el) == 1 {
		return el[0], "", nil
	}
	return el[0], el[1], nil
}
