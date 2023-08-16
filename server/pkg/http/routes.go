package http

import (
	"github.com/SergeyCherepiuk/chat-app/domain"
	"github.com/SergeyCherepiuk/chat-app/pkg/http/handlers"
	"github.com/SergeyCherepiuk/chat-app/pkg/http/middleware"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

type Router struct {
	AuthService          domain.AuthService
	UserService          domain.UserService
	DirectMessageService domain.DirectMessageService
	GroupChatService     domain.GroupChatService
}

func (router Router) Build() *fiber.App {
	app := fiber.New()

	api := app.Group("/api")

	authMiddleware := middleware.NewAuthMiddleware(router.AuthService)
	directMessageMiddleware := middleware.NewChatMiddleware(
		router.UserService,
		router.DirectMessageService,
	)
	groupChatMiddleware := middleware.NewGroupChatMiddleware(router.GroupChatService)

	authHandler := handlers.NewAuthHandler(router.AuthService)
	auth := api.Group("/auth")
	auth.Post("/signup", authHandler.SignUp)
	auth.Post("/login", authHandler.Login)
	auth.Post("/logout", authHandler.Logout)

	userHandler := handlers.NewUserHandler(router.UserService)
	user := api.Group("/user")
	user.Use(authMiddleware.CheckIfAuthenticated())
	user.Get("/me", userHandler.GetMe)
	user.Get("/:username", userHandler.GetByUsername)
	user.Put("/me", userHandler.UpdateMe)
	user.Delete("/me", userHandler.DeleteMe)

	directMessageHandler := handlers.NewDirectMessageHandler(
		router.DirectMessageService,
		router.UserService,
	)
	directChat := api.Group("/chat/:username")
	directChat.Use(authMiddleware.CheckIfAuthenticated())
	directChat.Use(directMessageMiddleware.CheckIfCompanionExists())
	directChat.Get(
		"/",
		middleware.Upgrade,
		websocket.New(directMessageHandler.EnterChat, websocket.Config{}),
	)
	directChat.Delete("/", directMessageHandler.DeleteChat)

	directMessage := directChat.Group("/:message_id")
	directMessage.Use(directMessageMiddleware.CheckIfBelongsToChat())
	directMessage.Put("/", directMessageMiddleware.CheckIfAuthor(), directMessageHandler.UpdateMessage)
	directMessage.Delete("/", directMessageHandler.DeleteMessage)

	groupChatHandler := handlers.NewGroupChatHandler(router.GroupChatService)
	groupChat := api.Group("/group-chat")
	groupChat.Use(authMiddleware.CheckIfAuthenticated())
	groupChat.Post("/", groupChatHandler.CreateChat)

	groupChat.Get(
		"/:chat_id/enter",
		middleware.Upgrade,
		groupChatMiddleware.CheckIfGroupChatExists(),
		websocket.New(groupChatHandler.EnterChat, websocket.Config{}),
	)

	groupChatWithId := groupChat.Group("/:chat_id")
	groupChatWithId.Use(groupChatMiddleware.CheckIfGroupChatExists())
	groupChatWithId.Get("/", groupChatHandler.GetChat)
	groupChatWithId.Put("/", groupChatMiddleware.CheckIfAdmin(), groupChatHandler.UpdateChat)
	groupChatWithId.Delete("/", groupChatMiddleware.CheckIfAdmin(), groupChatHandler.DeleteChat)

	groupMessage := groupChatWithId.Group("/:message_id")
	groupMessage.Use(groupChatMiddleware.CheckIfMessageBelongsToChat())
	groupMessage.Put("/", groupChatMiddleware.CheckIfAuthorOfMessage(), groupChatHandler.UpdateMessage)
	groupMessage.Delete("/", groupChatMiddleware.CheckIfAuthorOfMessage(), groupChatHandler.DeleteMessage)

	return app
}
