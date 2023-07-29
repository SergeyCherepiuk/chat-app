package authhandler_test

import (
	authhandler "github.com/SergeyCherepiuk/chat-app/handlers/auth"
	authstorage "github.com/SergeyCherepiuk/chat-app/storage/auth"
	"github.com/gofiber/fiber/v2"
)

var app *fiber.App

func init() {
	app = fiber.New()
	
	authStorageMock := authstorage.NewMock()
	authHandler := authhandler.New(authStorageMock)
	
	app.Post("/login", authHandler.Login)
	app.Post("/signup", authHandler.SignUp)
	app.Post("/logout", authHandler.Logout)
}
