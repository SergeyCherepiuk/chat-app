package handlers

import (
	"github.com/SergeyCherepiuk/session-auth/src/auth"
	"github.com/SergeyCherepiuk/session-auth/src/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserHandler struct {
	pdb            *gorm.DB
	sessionManager *auth.SessionManager
}

func NewUserHandler(pdb *gorm.DB, sessionManager *auth.SessionManager) *UserHandler {
	return &UserHandler{pdb: pdb, sessionManager: sessionManager}
}

type meResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
}

func (handler UserHandler) Me(c *fiber.Ctx) error {
	sessionId, err := uuid.Parse(c.Cookies("session_id", ""))
	if err != nil {
		return err
	}

	userId, err := handler.sessionManager.CheckSession(sessionId)
	if err != nil {
		return err
	}

	user := models.User{}
	handler.pdb.First(&user, userId)

	return c.JSON(meResponse{
		ID:       user.ID,
		Username: user.Username,
	})
}
