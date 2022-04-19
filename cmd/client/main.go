// Entry point for application's clients
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/streadway/amqp"
	"github.com/vbetsun/rmq-dynamic-clients/configs"
	"github.com/vbetsun/rmq-dynamic-clients/internal/net/rpc"
)

func main() {
	cfg, err := configs.New()
	if err != nil {
		log.Fatal(err)
	}
	conn, err := amqp.Dial(cfg.AMQP.Url)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	clientCodec, err := rpc.NewClientCodec(conn, cfg, rpc.GobCodec{})
	if err != nil {
		log.Fatal(err)
	}

	cmdRouter := rpc.BuildRouter()
	client := rpc.NewClientWithCodec(clientCodec)
	scanner := bufio.NewScanner(os.Stdin)
	log.Println("Client is running")
	for scanner.Scan() {
		cmd, arg, err := parseCmd(scanner.Text())
		if err != nil {
			log.Println(err)
			continue
		}
		if handler, ok := cmdRouter[cmd]; ok {
			err = handler(client, arg)
			if err != nil {
				log.Print(err)
			}
		} else {
			log.Printf("unknown command %s", cmd)
		}
	}
}

func parseCmd(cmd string) (string, string, error) {
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
