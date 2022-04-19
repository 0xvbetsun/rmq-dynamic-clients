# RMQ dynamic clients

![CI](https://github.com/vbetsun/rmq-dynamic-clients/workflows/CI/badge.svg)
[![GoReport](https://goreportcard.com/badge/github.com/vbetsun/rmq-dynamic-clients)](https://goreportcard.com/report/github.com/vbetsun/rmq-dynamic-clients)
![license](https://img.shields.io/github/license/vbetsun/rmq-dynamic-clients)
[![GoDoc](https://pkg.go.dev/badge/github.com/vbetsun/rmq-dynamic-clients)](https://pkg.go.dev/github.com/vbetsun/rmq-dynamic-clients)

## How to Deploy

```sh
cp .env.sample .env

docker-compose -f ./deployments/docker-compose.yml up
```

## How to Test

```sh
docker attach client

# write your command and press ENTER
AddItem abc
AddItem bcd
AddItem cde

GetItem bcd

GetAllItems
RemoveItem bcd
GetAllItems
```