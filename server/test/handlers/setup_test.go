package handlers_test

import (
	"github.com/SergeyCherepiuk/chat-app/mocks"
	"github.com/SergeyCherepiuk/chat-app/pkg/http"
	"github.com/SergeyCherepiuk/chat-app/pkg/log"
	"github.com/gofiber/fiber/v2"
)

var app *fiber.App

func init() {
	app = http.Router{
		AuthService:                    mocks.NewAuthService(),
		UserService:                    mocks.NewUserService(),
		DirectMessageService:           mocks.NewDirectMessageService(),
		DirectConnectionManagerService: mocks.NewConnectionManagerService[[2]uint](),
		GroupChatService:               mocks.NewGroupChatService(),
		GroupConnectionManagerService:  mocks.NewConnectionManagerService[uint](),
	}.Build()

	go log.HandleLogs()
}
