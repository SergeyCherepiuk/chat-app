FROM golang:latest

RUN apt-get update && apt-get install -y gcc

RUN go install github.com/pressly/goose/v3/cmd/goose@latest

RUN mkdir -p /opt/migrations
COPY ./migrations /opt/migrations
WORKDIR /opt/migrations

ENTRYPOINT [ "/go/bin/goose" ]
CMD [ "postgres", "user=postgres password=postgres host=postgres-db dbname=chat_app sslmode=disable", "up" ]