package http

import (
	"github.com/SergeyCherepiuk/chat-app/domain"
	"github.com/SergeyCherepiuk/chat-app/pkg/http/handlers"
	"github.com/SergeyCherepiuk/chat-app/pkg/http/middleware"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func NewRouter(
	sessionManagerService domain.SessionManagerService,
	authService domain.AuthService,
	directMessageService domain.DirectMessageService,
	groupChatService domain.GroupChatService,
	userService domain.UserService,
) *fiber.App {
	app := fiber.New()

	api := app.Group("/api")

	authMiddleware := middleware.NewAuthMiddleware(authService)
	chatMiddleware := middleware.NewChatMiddleware(userService, directMessageService)
	groupChatMiddleware := middleware.NewGroupChatMiddleware(groupChatService)

	authHandler := handlers.NewAuthHandler(authService)
	auth := api.Group("/auth")
	auth.Post("/signup", authHandler.SignUp)
	auth.Post("/login", authHandler.Login)
	auth.Post("/logout", authHandler.Logout)

	userHandler := handlers.NewUserHandler(userService)
	user := api.Group("/user")
	user.Use(authMiddleware.CheckIfAuthenticated())
	user.Get("/me", userHandler.GetMe)
	user.Get("/:username", userHandler.GetByUsername)
	user.Put("/me", userHandler.UpdateMe)
	user.Delete("/me", userHandler.DeleteMe)

	chatHandler := handlers.NewDirectMessageHandler(directMessageService, userService)
	chat := api.Group("/chat/:username")
	chat.Use(authMiddleware.CheckIfAuthenticated())
	chat.Use(chatMiddleware.CheckIfCompanionExists())
	chat.Delete("/", chatHandler.DeleteChat)

	message := chat.Group("/:message_id")
	message.Use(chatMiddleware.CheckIfBelongsToChat())
	message.Put("/", chatMiddleware.CheckIfAuthor(), chatHandler.UpdateMessage)
	message.Delete("/", chatHandler.DeleteMessage)

	ws := chat.Group("")
	ws.Use(middleware.Upgrade)
	ws.Get("/", websocket.New(chatHandler.EnterChat, websocket.Config{}))

	groupChatHandler := handlers.NewGroupChatHandler(groupChatService)
	groupChat := api.Group("/group-chat")
	groupChat.Use(authMiddleware.CheckIfAuthenticated())
	groupChat.Get("/:chat_id", groupChatHandler.GetChat)
	groupChat.Post("/", groupChatHandler.Create)
	groupChat.Put("/:chat_id", groupChatMiddleware.CheckIfAdmin(), groupChatHandler.Update)
	groupChat.Delete("/:chat_id", groupChatMiddleware.CheckIfAdmin(), groupChatHandler.Delete)

	return app
}
