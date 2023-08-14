package main

import (
	stdlog "log"

	"github.com/SergeyCherepiuk/chat-app/pkg/database/postgres"
	"github.com/SergeyCherepiuk/chat-app/pkg/database/redis"
	"github.com/SergeyCherepiuk/chat-app/pkg/http"
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
		AuthService:          postgres.NewAuthService(redis.NewSessionManagerService()),
		UserService:          postgres.NewUserService(),
		DirectMessageService: postgres.NewDirectMessageService(),
		GroupChatService:     postgres.NewGroupChatService(),
	}.Build()

	app.Listen(":8001")
}
