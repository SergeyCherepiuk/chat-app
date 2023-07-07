package main

import (
	"github.com/SergeyCherepiuk/session-auth/auth"
	"github.com/SergeyCherepiuk/session-auth/handlers"
	"github.com/SergeyCherepiuk/session-auth/initializers"
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

	authHandler := handlers.NewAuthHandler(pdb, sessionManager)
	auth := api.Group("/auth")
	auth.Post("/signup", authHandler.SingUp)
	auth.Post("/login", authHandler.Login)
	auth.Post("/logout", authHandler.Logout)

	userHandler := handlers.NewUserHandler(pdb, sessionManager)
	user := api.Group("/user")
	user.Get("/me", userHandler.Me)

	rolesHandler := handlers.NewRolesHandler(pdb)
	roles := api.Group("/roles")
	roles.Get("/:id", rolesHandler.GetById)
	roles.Get("/", rolesHandler.GetAll)
	roles.Post("/", rolesHandler.Create)
	roles.Put("/:id", rolesHandler.Update)
	roles.Delete("/:id", rolesHandler.Delete)
	roles.Delete("/", rolesHandler.DeleteAll)

	app.Listen(":8001")
}
