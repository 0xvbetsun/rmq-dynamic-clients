// Package configs provides project configuration structure
package configs

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
)

// Configuration for RabbitMQ
type AMQP struct {
	Url        string `json:"url"`
	Queue      string `json:"queue"`
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
	p := filepath.Join("configs", "conf.json")
	file, err := ioutil.ReadFile(filepath.Clean(p))
	if err != nil {
		return nil, err
	}
	conf := &Config{}
	if err = json.Unmarshal(file, conf); err != nil {
		return nil, err
	}

	amqpURL, amqpQueue, docsPort := os.Getenv("AMQP_SERVER_URL"), os.Getenv("AMQP_QUEUE_NAME"), os.Getenv("DOCS_PORT")
	if amqpURL != "" {
		conf.AMQP.Url = amqpURL
	}
	if amqpQueue != "" {
		conf.AMQP.Queue = amqpQueue
	}

	if docsPort != "" {
		if conf.Docs.Port, err = strconv.Atoi(docsPort); err != nil {
			return nil, err
		}
	}

	return conf, nil
}
