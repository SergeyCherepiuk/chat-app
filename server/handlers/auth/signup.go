package authhandler

import (
	"errors"
	"strings"
	"time"

	"github.com/SergeyCherepiuk/chat-app/logger"
	"github.com/SergeyCherepiuk/chat-app/models"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/exp/slog"
)

type SignUpRequestBody struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
	Password  string `json:"password"`
}

func (body SignUpRequestBody) Validate() error {
	var err error

	if strings.TrimSpace(body.FirstName) == "" {
		err = errors.Join(err, errors.New("first name is empty"))
	}

	if strings.TrimSpace(body.LastName) == "" {
		err = errors.Join(err, errors.New("last name is empty"))
	}

	if strings.TrimSpace(body.Username) == "" {
		err = errors.Join(err, errors.New("username is empty"))
	}

	if strings.TrimSpace(body.Password) == "" {
		err = errors.Join(err, errors.New("password is empty"))
	} else if len(body.Password) < 8 {
		err = errors.Join(err, errors.New("password is too short"))
	}

	return err
}

func (handler AuthHandler) SignUp(c *fiber.Ctx) error {
	body := SignUpRequestBody{}
	if err := c.BodyParser(&body); err != nil {
		logger.LogMessages <- logger.LogMessage{
			Message: "failed to parse the body",
			Level:   slog.LevelError,
			Attrs:   []slog.Attr{slog.String("err", err.Error())},
		}
		return err
	}

	if err := body.Validate(); err != nil {
		logger.LogMessages <- logger.LogMessage{
			Message: "request body isn't valid",
			Level:   slog.LevelError,
			Attrs: []slog.Attr{
				slog.String("err", err.Error()),
				slog.Any("body", body),
			},
		}
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		logger.LogMessages <- logger.LogMessage{
			Message: "failed to hash the password",
			Level:   slog.LevelError,
			Attrs: []slog.Attr{
				slog.String("err", err.Error()),
				slog.String("password", body.Password),
			},
		}
		return err
	}

	user := models.User{
		FirstName: body.FirstName,
		LastName:  body.LastName,
		Username:  body.Username,
		Password:  string(hashedPassword),
	}
	sessionId, err := handler.storage.SignUp(&user)
	if err != nil {
		logger.LogMessages <- logger.LogMessage{
			Message: "failed to sign up the user",
			Level:   slog.LevelError,
			Attrs: []slog.Attr{
				slog.String("err", err.Error()),
				slog.Any("user", user),
			},
		}
		return err
	}

	logger.LogMessages <- logger.LogMessage{
		Message: "user has been signed up",
		Level:   slog.LevelInfo,
		Attrs: []slog.Attr{
			slog.Uint64("user_id", uint64(user.ID)),
			slog.Any("session_id", sessionId),
		},
	}
	c.Cookie(&fiber.Cookie{
		Name:     "session_id",
		Value:    sessionId.String(),
		HTTPOnly: true,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
	})
	return c.SendStatus(fiber.StatusOK)
}
