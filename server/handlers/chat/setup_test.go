package chathandler_test

import (
	chathandler "github.com/SergeyCherepiuk/chat-app/handlers/chat"
	"github.com/SergeyCherepiuk/chat-app/logger"
	"github.com/SergeyCherepiuk/chat-app/middleware"
	authstorage "github.com/SergeyCherepiuk/chat-app/storage/auth"
	chatstorage "github.com/SergeyCherepiuk/chat-app/storage/chat"
	userstorage "github.com/SergeyCherepiuk/chat-app/storage/user"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

var app *fiber.App

func init() {
	app = fiber.New()

	authStorage := authstorage.NewMock()
	chatStorage := chatstorage.NewMock()
	userStorage := userstorage.NewMock()

	authMiddleware := middleware.NewAuthMiddleware(authStorage)
	chatMiddleware := middleware.NewChatMiddleware(userStorage, chatStorage)

	chatHandler := chathandler.New(chatStorage, userStorage)
	chat := app.Group("/chat/:username")
	chat.Use(authMiddleware.CheckIfAuthenticated())
	chat.Use(chatMiddleware.CheckIfCompanionExists())
	chat.Delete("/", chatHandler.DeleteChat)

	message := chat.Group("/:message_id")
	message.Use(chatMiddleware.CheckIfBelongsToChat())
	message.Put("/", chatMiddleware.CheckIfAuthor(), chatHandler.UpdateMessage)
	message.Delete("/", chatHandler.DeleteMessage)

	ws := chat.Group("")
	ws.Use(middleware.Upgrade)
	ws.Get("/", websocket.New(chatHandler.EnterChat, websocket.Config{}))

	go logger.HandleLogs()
}