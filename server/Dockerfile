FROM golang:1.20.4-alpine

WORKDIR /server

COPY go.mod go.sum ./
RUN go mod download

COPY . .
COPY "../.env" .

CMD [ "go", "run", "main.go" ]
EXPOSE 8001