package main

import (
	"github.com/SergeyCherepiuk/session-auth/src/handlers"
	"github.com/SergeyCherepiuk/session-auth/src/initializers"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"gorm.io/gorm"
)

var db *gorm.DB

func init() {
	initializers.LoadEnv()
	db = initializers.MustConnect()
}

func main() {
	app := fiber.New()
	app.Use(logger.New())

	api := app.Group("/api")
	user := api.Group("/user")

	authHandler := handlers.NewAuthHandler(db)
	auth := api.Group("/auth")
	auth.Post("/signup", authHandler.SingUp)
	auth.Post("/login", authHandler.Login)
	auth.Post("/logout", authHandler.Logout)

	userHandler := handlers.NewUserHandler(db)
	user.Get("/me", userHandler.Me)

	app.Listen(":8001")
}
