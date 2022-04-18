package main

import (
	"bufio"
	"fmt"
	"log"
	"net/rpc"
	"os"
	"strings"

	"github.com/streadway/amqp"
	"github.com/vbetsun/rmq-dynamic-clients/internal/codec"
)

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
	clientCodec, err := codec.NewClientCodec(conn, amqpQueue, codec.GobCodec{})
	if err != nil {
		log.Fatal(err)
	}

	cmdRouter := buildRouter()
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

func AddItem(c *rpc.Client, arg string) error {
	return c.Call("Items.AddItem", arg, nil)
}

func GetItem(c *rpc.Client, arg string) error {
	return c.Call("Items.GetItem", arg, nil)
}

func GetAllItems(c *rpc.Client, arg string) error {
	return c.Call("Items.GetAllItems", arg, nil)
}

func RemoveItem(c *rpc.Client, arg string) error {
	return c.Call("Items.RemoveItem", arg, nil)
}

func buildRouter() map[string]func(*rpc.Client, string) error {
	return map[string]func(*rpc.Client, string) error{
		"AddItem":     AddItem,
		"RemoveItem":  RemoveItem,
		"GetItem":     GetItem,
		"GetAllItems": GetAllItems,
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
