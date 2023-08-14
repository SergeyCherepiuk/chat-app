package main

import (
	stdlog "log"

	"github.com/SergeyCherepiuk/chat-app/pkg/connection"
	"github.com/SergeyCherepiuk/chat-app/pkg/database/postgres"
	"github.com/SergeyCherepiuk/chat-app/pkg/database/redis"
	"github.com/SergeyCherepiuk/chat-app/pkg/http"
	"github.com/SergeyCherepiuk/chat-app/pkg/log"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		stdlog.Fatal(err)
	}

	postgres.PostgresMustConnect()
	redis.RedisMustConnect()
}

func main() {
	app := http.Router{
		AuthService: postgres.NewAuthService(
			redis.NewSessionManagerService(),
		),
		UserService: postgres.NewUserService(),
		DirectMessageService: postgres.NewDirectMessageService(),
		DirectConnectionManagerService: connection.NewConnectionManager[[2]uint](),
		GroupChatService: postgres.NewGroupChatService(),
		GroupConnectionManagerService: connection.NewConnectionManager[uint](),
	}.Build()

	for i := 0; i < 10; i++ {
		go log.HandleLogs()
	}

	app.Listen(":8001")
}
