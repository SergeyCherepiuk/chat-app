FROM golang:latest

RUN apt-get update && apt-get install -y gcc

RUN go install github.com/pressly/goose/v3/cmd/goose@latest

RUN mkdir -p /opt/migrations
COPY ./migrations /opt/migrations
WORKDIR /opt/migrations

ARG POSTGRES_USER
ARG POSTGRES_PASSWORD
ARG POSTGRES_DBNAME

ENV DBSTRING="user=$POSTGRES_USER password=$POSTGRES_PASSWORD host=postgres-db dbname=$POSTGRES_DBNAME sslmode=disable"

CMD [ "sh", "-c", "/go/bin/goose postgres \"$DBSTRING\" up" ]