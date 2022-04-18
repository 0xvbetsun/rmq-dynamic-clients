# RMQ dynamic clients

## How to Deploy

```sh
cp ./configs/.env.sample ./configs/.env

docker-compose --env-file ./configs/.env -f ./deployments/docker-compose.yml up
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