package main

import (
	authhandler "github.com/SergeyCherepiuk/chat-app/handlers/auth"
	chathandler "github.com/SergeyCherepiuk/chat-app/handlers/chat"
	userhandler "github.com/SergeyCherepiuk/chat-app/handlers/user"
	"github.com/SergeyCherepiuk/chat-app/initializers"
	"github.com/SergeyCherepiuk/chat-app/logger"
	"github.com/SergeyCherepiuk/chat-app/middleware"
	authstorage "github.com/SergeyCherepiuk/chat-app/storage/auth"
	chatstorage "github.com/SergeyCherepiuk/chat-app/storage/chat"
	userstorage "github.com/SergeyCherepiuk/chat-app/storage/user"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

var rdb *redis.Client
var pdb *sqlx.DB

func init() {
	initializers.LoadEnv()
	rdb = initializers.RedisMustConnect()
	pdb = initializers.PostgresMustConnect()
}

func main() {
	app := fiber.New()

	api := app.Group("/api")

	sessionManager := authstorage.NewSessionManager(rdb)
	authStorage := authstorage.New(pdb, sessionManager)
	chatStorage := chatstorage.New(pdb)
	userStorage := userstorage.New(pdb, rdb)

	authMiddleware := middleware.NewAuthMiddleware(authStorage)
	chatMiddleware := middleware.NewChatMiddleware(userStorage, chatStorage)

	authHandler := authhandler.New(authStorage)
	auth := api.Group("/auth")
	auth.Post("/signup", authHandler.SignUp)
	auth.Post("/login", authHandler.Login)
	auth.Post("/logout", authHandler.Logout)

	userHandler := userhandler.New(userStorage)
	user := api.Group("/user")
	user.Use(authMiddleware.CheckIfAuthenticated())
	user.Get("/me", userHandler.GetMe)
	user.Get("/:username", userHandler.GetByUsername)
	user.Put("/me", userHandler.UpdateMe)
	user.Delete("/me", userHandler.DeleteMe)

	chatHandler := chathandler.New(chatStorage, userStorage)
	chat := api.Group("/chat/:username")
	chat.Use(authMiddleware.CheckIfAuthenticated())
	chat.Use(chatMiddleware.CheckIfCompanionExists())
	chat.Delete("/", chatHandler.DeleteChat)

	message := chat.Group("/:message_id")
	message.Use(chatMiddleware.CheckIfBelongsToChat())
	message.Put("/", chatHandler.UpdateMessage)
	message.Delete("/", chatHandler.DeleteMessage)

	ws := chat.Group("")
	ws.Use(middleware.Upgrade)
	ws.Get("/", websocket.New(chatHandler.EnterChat, websocket.Config{}))

	for i := 0; i < 10; i++ {
		go logger.HandleLogs()
	}

	app.Listen(":8001")
}
