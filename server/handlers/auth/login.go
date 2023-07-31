package authhandler

import (
	"errors"
	"strings"
	"time"

	"github.com/SergeyCherepiuk/chat-app/logger"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/exp/slog"
)

type LoginRequestBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (body LoginRequestBody) Validate() error {
	var err error

	if strings.TrimSpace(body.Username) == "" {
		err = errors.Join(errors.New("username is empty"))
	}

	if strings.TrimSpace(body.Password) == "" {
		err = errors.Join(err, errors.New("password is empty"))
	} else if len(body.Password) < 8 {
		err = errors.Join(err, errors.New("password is too short"))
	}

	return err
}

func (handler AuthHandler) Login(c *fiber.Ctx) error {
	log := logger.Logger{}

	body := LoginRequestBody{}
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

	sessionId, userId, err := handler.storage.Login(body.Username, body.Password)
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
