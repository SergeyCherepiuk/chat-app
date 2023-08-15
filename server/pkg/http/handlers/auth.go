package handlers

import (
	"time"

	"github.com/SergeyCherepiuk/chat-app/domain"
	"github.com/SergeyCherepiuk/chat-app/pkg/http/validation"
	"github.com/SergeyCherepiuk/chat-app/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/exp/slog"
)

type AuthHandler struct {
	authService domain.AuthService
}

func NewAuthHandler(authService domain.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (handler AuthHandler) SignUp(c *fiber.Ctx) error {
	log := logger.Logger{}

	body := validation.SignUpRequestBody{}
	if err := c.BodyParser(&body); err != nil {
		log.Error("failed to parse the body", slog.String("err", err.Error()))
		return err
	}

	if err := body.Validate(); err != nil {
		log.Error(
			"request body isn't valid",
			slog.String("err", err.Error()),
			slog.Any("body", body),
		)
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		log.Error(
			"failed to hash the password",
			slog.String("err", err.Error()),
			slog.String("password", body.Password),
		)
		return err
	}

	user := domain.User{
		FirstName:      body.FirstName,
		LastName:       body.LastName,
		Username:       body.Username,
		Password:       string(hashedPassword),
		Description:    body.Description,
		ProfilePicture: body.ProfilePicture,
	}
	sessionId, userId, err := handler.authService.SignUp(user)
	if err != nil {
		log.Error(
			"failed to sign up the user",
			slog.String("err", err.Error()),
			slog.Any("user", user),
		)
		return err
	}

	log.Info(
		"user has been signed up",
		slog.Uint64("user_id", uint64(userId)),
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
	log := logger.Logger{}

	body := validation.LoginRequestBody{}
	if err := c.BodyParser(&body); err != nil {
		log.Error("failed to parse the body", slog.String("err", err.Error()))
		return err
	}

	sessionId, userId, err := handler.authService.Login(body.Username, body.Password)
	if err != nil {
		log.Error("failed to log in user", slog.String("err", err.Error()))
		return err
	}

	log.Info(
		"user has been logged in",
		slog.Uint64("user_id", uint64(userId)),
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
	log := logger.Logger{}

	sessionId, err := uuid.Parse(c.Cookies("session_id", ""))
	if err != nil {
		log.Error(
			"invalid session id",
			slog.String("err", err.Error()),
			slog.Any("session_id", sessionId),
		)
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	if err := handler.authService.Logout(sessionId); err != nil {
		log.Error("failed to log out user", slog.String("err", err.Error()))
		return err
	}

	log.Info("user has been logged out", slog.Any("session_id", sessionId))
	c.Cookie(&fiber.Cookie{
		Name:    "session_id",
		Expires: time.Now(),
	})
	return c.SendStatus(fiber.StatusOK)
}
