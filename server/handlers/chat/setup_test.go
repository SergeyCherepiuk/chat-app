package chathandler_test

import (
	chathandler "github.com/SergeyCherepiuk/chat-app/handlers/chat"
	"github.com/SergeyCherepiuk/chat-app/logger"
	"github.com/SergeyCherepiuk/chat-app/middleware"
	authstorage "github.com/SergeyCherepiuk/chat-app/storage/auth"
	chatstorage "github.com/SergeyCherepiuk/chat-app/storage/chat"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

var app *fiber.App

func init() {
	app = fiber.New()

	authStorageMock := authstorage.NewMock()
	authMiddleware := middleware.NewAuthMiddleware(authStorageMock)
	chatStorageMock := chatstorage.NewMock()
	chatHandler := chathandler.New(chatStorageMock)

	app.Use(authMiddleware.CheckIfAuthenticated())
	app.Get("/", chatHandler.GetAll)
	app.Get("/:chat_id", chatHandler.GetById)
	app.Post("/", chatHandler.Create)
	app.Put("/:chat_id", chatHandler.Update)
	app.Delete("/:chat_id", chatHandler.Delete)

	ws := app.Group("")
	ws.Use(middleware.Upgrade)
	ws.Get("/:chat_id/enter", websocket.New(chatHandler.Enter, websocket.Config{}))

	go logger.HandleLogs()
}
