package main

import (
	"github.com/SergeyCherepiuk/chat-app/authentication/handlers"
	"github.com/SergeyCherepiuk/chat-app/authentication/initializers"
	"github.com/SergeyCherepiuk/chat-app/authentication/middleware"
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
	authHandler := handlers.NewAuthHandler(authStorage)
	auth := api.Group("/auth")
	auth.Post("/signup", authHandler.SignUp)
	auth.Post("/login", authHandler.Login)
	auth.Post("/logout", authHandler.Logout)

	authMiddleware := middleware.NewAuthMiddleware(authStorage)
	userStorage := storage.NewUserStorage(pdb)
	userHandler := handlers.NewUserHandler(userStorage)
	user := api.Group("/user")
	user.Use(authMiddleware.CheckIfAuthenticated())
	user.Get("/me", userHandler.GetMe)
	user.Get("/:username", userHandler.GetUser)
	user.Put("/me", userHandler.UpdateMe)
	user.Delete("/me", userHandler.DeleteMe)

	app.Listen(":8001")
}
