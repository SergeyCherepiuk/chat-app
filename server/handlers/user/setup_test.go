package userhandler_test

import (
	userhandler "github.com/SergeyCherepiuk/chat-app/handlers/user"
	"github.com/SergeyCherepiuk/chat-app/middleware"
	authstorage "github.com/SergeyCherepiuk/chat-app/storage/auth"
	userstorage "github.com/SergeyCherepiuk/chat-app/storage/user"
	"github.com/gofiber/fiber/v2"
)

var app *fiber.App

func init() {
	app = fiber.New()

	authStorageMock := authstorage.NewMock()
	authMiddleware := middleware.NewAuthMiddleware(authStorageMock)
	userStorageMock := userstorage.NewMock()
	userHandler := userhandler.New(userStorageMock)

	app.Use(authMiddleware.CheckIfAuthenticated())
	app.Get("/me", userHandler.GetMe)
	app.Get("/:username", userHandler.GetByUsername)
	app.Put("/me", userHandler.UpdateMe)
	app.Delete("/me", userHandler.DeleteMe)
}