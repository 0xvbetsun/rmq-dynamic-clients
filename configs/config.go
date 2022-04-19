// Package configs provides project configuration structure
package configs

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"strconv"
)

// Configuration for RabbitMQ
type AMQP struct {
	Url        string `json:"-"`
	Queue      string `json:"-"`
	Durable    bool   `json:"durable"`
	AutoDelete bool   `json:"autoDelete"`
	AutoAck    bool   `json:"autoAck"`
	Exclusive  bool   `json:"exclusive"`
	NoWait     bool   `json:"noWait"`
	NoLocal    bool   `json:"noLocal"`
}

// Configuration for AsyncApi documentation
type Docs struct {
	Port int `json:"port"`
}

type Config struct {
	AMQP *AMQP `json:"amqp"`
	Docs *Docs `json:"docs"`
}

// New reads and validates env variables, config file and returns new instance of Config
func New() (*Config, error) {
	file, err := ioutil.ReadFile("./configs/conf.json")
	if err != nil {
		return nil, err
	}
	conf := &Config{}
	if err = json.Unmarshal(file, conf); err != nil {
		return nil, err
	}

	amqpURL, amqpQueue, docsPort := os.Getenv("AMQP_SERVER_URL"), os.Getenv("AMQP_QUEUE_NAME"), os.Getenv("DOCS_PORT")
	if amqpURL == "" {
		return nil, errors.New("env var AMQP_SERVER_URL was not provided")
	}
	conf.AMQP.Url = amqpURL
	if amqpQueue == "" {
		return nil, errors.New("env var AMQP_QUEUE_NAME was not provided")
	}
	conf.AMQP.Queue = amqpQueue
	if docsPort == "" {
		docsPort = "3000"
	}
	if conf.Docs.Port, err = strconv.Atoi(docsPort); err != nil {
		return nil, err
	}
	return conf, nil
}
