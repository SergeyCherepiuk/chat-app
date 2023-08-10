package handlers_test

import (
	"github.com/SergeyCherepiuk/chat-app/mocks"
	"github.com/SergeyCherepiuk/chat-app/pkg/http"
	"github.com/SergeyCherepiuk/chat-app/pkg/log"
	"github.com/gofiber/fiber/v2"
)

var app *fiber.App

func init() {
	app = http.NewRouter(
		mocks.NewSessionManagerService(),
		mocks.NewAuthService(),
		mocks.NewDirectMessageService(),
		mocks.NewGroupChatService(),
		mocks.NewUserService(),
	)

	go log.HandleLogs()
}
