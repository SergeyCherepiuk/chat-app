# Session authentication

This repository is a pet project aimed at exploring microservice architecture and session management.

## Installation and usage

Run with Docker Compose using the command below

```shell
docker-compose up --build
```

or

1. Clone the repository

2. Make sure your PostgreSQL and Redis databases are running

3. Porvide the way to connect to your databases in `.env` file (you can find example in `example.env`)

4. Run the server
```shell
cd server
go run main.go
```

## Check endpoints in Postman

[Postman collection](https://www.postman.com/asdkoda/workspace/session-auth/collection/26414129-53161e41-e3df-473f-81c7-4be929c6500e?action=share&creator=26414129)