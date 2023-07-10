package handlers

import (
	"time"

	"github.com/SergeyCherepiuk/chat-app/authentication/domain"
	"github.com/SergeyCherepiuk/chat-app/authentication/models"
	"github.com/SergeyCherepiuk/chat-app/authentication/storage"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
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
		return err
	}

	if err := body.Validate(); err != nil {
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
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
		return err
	}

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
		return err
	}

	if err := body.Validate(); err != nil {
		return err
	}

	sessionId, err := handler.storage.Login(body.Username, body.Password)
	if err != nil {
		return err
	}

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
		return err
	}

	if err := handler.storage.Logout(sessionId); err != nil {
		return err
	}

	c.Cookie(&fiber.Cookie{
		Name:    "session_id",
		Expires: time.Now(),
	})
	return c.SendStatus(fiber.StatusOK)
}
