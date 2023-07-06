package main

import (
	"github.com/SergeyCherepiuk/session-auth/src/auth"
	"github.com/SergeyCherepiuk/session-auth/src/handlers"
	"github.com/SergeyCherepiuk/session-auth/src/initializers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var pdb *gorm.DB
var rdb *redis.Client
var sessionManager *auth.SessionManager

func init() {
	initializers.LoadEnv()
	pdb = initializers.PostgresMustConnect()
	rdb = initializers.RedisMustConnect()
	sessionManager = auth.NewSessionManager(rdb)
}

func main() {
	app := fiber.New()
	app.Use(logger.New())

	api := app.Group("/api")
	user := api.Group("/user")

	authHandler := handlers.NewAuthHandler(pdb, rdb, sessionManager)
	auth := api.Group("/auth")
	auth.Post("/signup", authHandler.SingUp)
	auth.Post("/login", authHandler.Login)
	auth.Post("/logout", authHandler.Logout)

	userHandler := handlers.NewUserHandler(pdb, sessionManager)
	user.Get("/me", userHandler.Me)

	app.Listen(":8001")
}
