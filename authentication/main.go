package main

import (
	"github.com/SergeyCherepiuk/chat-app/authentication/handlers"
	"github.com/SergeyCherepiuk/chat-app/authentication/initializers"
	"github.com/SergeyCherepiuk/chat-app/authentication/storage"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
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
	app.Use(logger.New())

	api := app.Group("/api")
	
	authStorage := storage.NewAuthStorage(pdb, rdb)
	authHandler := handlers.NewAuthHandler(pdb, rdb, authStorage)
	auth := api.Group("/auth")
	auth.Post("/signup", authHandler.SignUp)
	auth.Post("/login", authHandler.Login)
	auth.Get("/check", authHandler.Check)
	auth.Post("/logout", authHandler.Logout)

	app.Listen(":8001")
}
