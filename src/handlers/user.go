package handlers

import (
	"errors"
	"strconv"
	"time"

	"github.com/SergeyCherepiuk/session-auth/src/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type UserHandler struct {
	db *gorm.DB
}

func NewUserHandler(db *gorm.DB) UserHandler {
	return UserHandler{db: db}
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

	session := models.Session{}
	if r := handler.db.First(&session, sessionId); r.Error != nil || r.RowsAffected < 1 {
		return errors.New("session not found")
	}

	if session.ExpiresAt.Before(time.Now().UTC()) {
		return errors.New("session expired")
	}

	user := models.User{}
	handler.db.First(&user, session.UserID)

	return c.JSON(meResponse{
		ID:       user.ID,
		Username: user.Username,
	})
}
