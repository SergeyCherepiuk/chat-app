FROM golang:1.20.4-alpine

WORKDIR /authentication

COPY go.mod go.sum ./
RUN go mod download

COPY . .

CMD [ "go", "run", "main.go" ]