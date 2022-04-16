# RMQ dynamic clients

## How to DEploy

```sh
cp ./configs/.env.sample ./configs/.env

docker-compose --env-file ./configs/.env -f ./deployments/docker-compose.yml up
```