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
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var rdb *redis.Client
var pdb *gorm.DB

func init() {
	initializers.LoadEnv()
	rdb = initializers.RedisMustConnect()
	pdb = initializers.PostgresMustConnect()
}

func main() {
	app := fiber.New()

	api := app.Group("/api")

	authStorage := authstorage.New(pdb, rdb)
	authHandler := authhandler.New(authStorage)
	auth := api.Group("/auth")
	auth.Post("/signup", authHandler.SignUp)
	auth.Post("/login", authHandler.Login)
	auth.Post("/logout", authHandler.Logout)

	authMiddleware := middleware.NewAuthMiddleware(authStorage)
	userStorage := userstorage.New(pdb)
	userHandler := userhandler.New(userStorage)
	user := api.Group("/user")
	user.Use(authMiddleware.CheckIfAuthenticated())
	user.Get("/me", userHandler.GetMe)
	user.Get("/:username", userHandler.GetByUsername)
	user.Put("/me", userHandler.UpdateMe)
	user.Delete("/me", userHandler.DeleteMe)

	chat := api.Group("/chat")
	chatStorage := chatstorage.New(pdb)
	chatHandler := chathandler.New(chatStorage)
	chat.Use(authMiddleware.CheckIfAuthenticated())
	chat.Get("/", chatHandler.GetAll)
	chat.Get("/:chat_id", chatHandler.GetById)
	chat.Post("/", chatHandler.Create)
	chat.Put("/:chat_id", chatHandler.Update)
	chat.Delete("/:chat_id", chatHandler.Delete)

	ws := chat.Group("")
	ws.Use(middleware.Upgrade)
	ws.Get("/:chat_id/enter", websocket.New(chatHandler.Enter, websocket.Config{}))

	go logger.HandleLogs()

	app.Listen(":8001")
}
