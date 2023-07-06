package handlers

import (
	"strconv"

	"github.com/SergeyCherepiuk/session-auth/src/auth"
	"github.com/SergeyCherepiuk/session-auth/src/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type UserHandler struct {
	db             *gorm.DB
	sessionManager *auth.SessionManager
}

func NewUserHandler(db *gorm.DB, sessionManager *auth.SessionManager) *UserHandler {
	return &UserHandler{db: db, sessionManager: sessionManager}
}

type meResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
}

func (handler UserHandler) Me(c *fiber.Ctx) error {
	sessionId, err := strconv.ParseUint(c.Cookies("session_id", ""), 10, 64)
	if err != nil {
		return err
	}

	session, err := handler.sessionManager.CheckSession(uint(sessionId))
	if err != nil {
		return err
	}

	user := models.User{}
	handler.db.First(&user, session.UserID)

	return c.JSON(meResponse{
		ID:       user.ID,
		Username: user.Username,
	})
}
