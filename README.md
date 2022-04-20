# RMQ dynamic clients

The project represents the configuration and setup server with multiple (scalable) clients which communicate using RPC (Gob codec) over AMQP

![CI](https://github.com/vbetsun/rmq-dynamic-clients/workflows/CI/badge.svg)
[![GoReport](https://goreportcard.com/badge/github.com/vbetsun/rmq-dynamic-clients)](https://goreportcard.com/report/github.com/vbetsun/rmq-dynamic-clients)
![license](https://img.shields.io/github/license/vbetsun/rmq-dynamic-clients)
[![GoDoc](https://pkg.go.dev/badge/github.com/vbetsun/rmq-dynamic-clients)](https://pkg.go.dev/github.com/vbetsun/rmq-dynamic-clients)


## Prerequisites

- Git
- Docker
- Docker Compose
  
## How to Install

Clone from github and navigate to the project's folder
```sh
# HTTPS
git clone https://github.com/vbetsun/rmq-dynamic-clients.git

# SSH
git clone git@github.com:vbetsun/rmq-dynamic-clients.git

cd rmq-dynamic-clients
```

## How to Deploy

```sh
cp .env.sample .env
```

change env variables for your needs

```dotenv
DOCS_PORT=8080 # port for serving rpc documentation
AMQP_SERVER_URL=amqp://guest:guest@message-broker:5672 # RabbitMQ url
AMQP_QUEUE_NAME=queue_name # name of queue for usage
```

and start the application via docker-compose. It should start 1 instance of the server, 3 replicas of the client, and 1 instance of documentation service, which you can see on http://localhost:${DOCS_PORT}

```sh
docker-compose --env-file .env -f ./deployments/docker-compose.yml up -d
```
## How to Test

first of all, you have to run next command for making sure that everything is ok and status is "UP"

```sh
docker-compose --env-file .env -f ./deployments/docker-compose.yml ps -a
```

in the printed list you will see few clients (eg deployments_client_1, deployments_client_2, deployments_client_3)

```sh
docker attach deployments_client_1

# write your command and press ENTER
AddItem abc
AddItem bcd
AddItem cde

GetItem bcd
GetItem qwe

GetAllItems
RemoveItem bcd
GetAllItems

docker attach deployments_client_2
# and repeat previous step

docker attach deployments_client_3
# and repeat previous step
```

## How to scale

if you want to test the application with more than 3 replicas (for example 10), you can use the next command:

```sh
docker-compose --env-file .env -f ./deployments/docker-compose.yml up -d --scale client=10
# and repeat testing phase with attaching to the client_NUM
docker attach deployments_client_9
# and repeat testing
```