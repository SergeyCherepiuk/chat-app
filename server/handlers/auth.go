package handlers

import (
	"time"

	"github.com/SergeyCherepiuk/chat-app/domain"
	"github.com/SergeyCherepiuk/chat-app/logger"
	"github.com/SergeyCherepiuk/chat-app/models"
	"github.com/SergeyCherepiuk/chat-app/storage"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/exp/slog"
)

type AuthHandler struct {
	storage *storage.AuthStorage
}

func NewAuthHandler(storage *storage.AuthStorage) *AuthHandler {
	return &AuthHandler{storage: storage}
}

func (handler AuthHandler) SignUp(c *fiber.Ctx) error {
	body := domain.SignUpRequestBody{}
	if err := c.BodyParser(&body); err != nil {
		logger.Logger.Error(
			"failed to parse the body",
			slog.String("error_message", err.Error()),
		)
		return err
	}

	if err := body.Validate(); err != nil {
		logger.Logger.Error(
			"request body isn't valid",
			slog.String("error_message", err.Error()),
			slog.Any("body", body),
		)
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		logger.Logger.Error(
			"failed to hash the password",
			slog.String("error_message", err.Error()),
		)
		return err
	}

	user := models.User{
		FirstName: body.FirstName,
		LastName:  body.LastName,
		Username:  body.Username,
		Password:  string(hashedPassword),
	}
	sessionId, err := handler.storage.SignUp(user)
	if err != nil {
		logger.Logger.Error(
			"failed to sign up the user",
			slog.String("error_message", err.Error()),
		)
		return err
	}
	logger.Logger.Info(
		"user has been signed up",
		slog.Uint64("user_id", uint64(user.ID)),
		slog.Any("session_id", sessionId),
	)

	c.Cookie(&fiber.Cookie{
		Name:     "session_id",
		Value:    sessionId.String(),
		HTTPOnly: true,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
	})
	return c.SendStatus(fiber.StatusOK)
}

func (handler AuthHandler) Login(c *fiber.Ctx) error {
	body := domain.LoginRequestBody{}
	if err := c.BodyParser(&body); err != nil {
		logger.Logger.Error(
			"failed to hash the password",
			slog.String("error_message", err.Error()),
		)
		return err
	}

	if err := body.Validate(); err != nil {
		logger.Logger.Error(
			"request body isn't valid",
			slog.String("error_message", err.Error()),
			slog.Any("body", body),
		)
		return err
	}

	sessionId, err := handler.storage.Login(body.Username, body.Password)
	if err != nil {
		logger.Logger.Error(
			"failed to log in user",
			slog.String("error_message", err.Error()),
		)
		return err
	}
	logger.Logger.Info(
		"user has been logged in",
		slog.Any("session_id", sessionId),
	)

	c.Cookie(&fiber.Cookie{
		Name:     "session_id",
		Value:    sessionId.String(),
		HTTPOnly: true,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
	})
	return c.SendStatus(fiber.StatusOK)
}

func (handler AuthHandler) Logout(c *fiber.Ctx) error {
	sessionId, err := uuid.Parse(c.Cookies("session_id", ""))
	if err != nil {
		logger.Logger.Error(
			"invalid session id",
			slog.String("error_message", err.Error()),
			slog.Any("session_id", sessionId),
		)
		return err
	}

	if err := handler.storage.Logout(sessionId); err != nil {
		logger.Logger.Error(
			"failed to log out user",
			slog.String("error_message", err.Error()),
		)
		return err
	}
	logger.Logger.Info("user has been logged out")

	c.Cookie(&fiber.Cookie{
		Name:    "session_id",
		Expires: time.Now(),
	})
	return c.SendStatus(fiber.StatusOK)
}
