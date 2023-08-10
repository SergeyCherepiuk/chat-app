package main

import (
	stdlog "log"

	"github.com/SergeyCherepiuk/chat-app/database/postgres"
	"github.com/SergeyCherepiuk/chat-app/database/redis"
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
	sessionManagerService := redis.NewSessionManagerService()

	app := http.NewRouter(
		sessionManagerService,
		postgres.NewAuthService(sessionManagerService),
		postgres.NewDirectMessageService(),
		postgres.NewGroupChatService(),
		postgres.NewUserService(),
	)

	for i := 0; i < 10; i++ {
		go log.HandleLogs()
	}

	app.Listen(":8001")
}
