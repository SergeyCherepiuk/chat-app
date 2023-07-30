package chathandler

import (
	"errors"
	"strings"

	"github.com/SergeyCherepiuk/chat-app/logger"
	"github.com/SergeyCherepiuk/chat-app/models"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/exp/slog"
)

type CreateChatRequestBody struct {
	Name string `json:"name"`
}

func (body CreateChatRequestBody) Validate() error {
	var err error

	if strings.TrimSpace(body.Name) == "" {
		err = errors.Join(err, errors.New("name is empty"))
	}

	return err
}

func (handler ChatHandler) Create(c *fiber.Ctx) error {
	userId, _ := c.Locals("user_id").(uint)

	body := CreateChatRequestBody{}
	if err := c.BodyParser(&body); err != nil {
		logger.LogMessages <- logger.LogMessage{
			Message: "failed to parse request body",
			Level:   slog.LevelError,
			Attrs: []slog.Attr{
				slog.String("err", err.Error()),
				slog.Uint64("user_id", uint64(userId)),
			},
		}
		return err
	}

	chat := models.Chat{Name: body.Name}
	if err := handler.storage.CreateChat(&chat); err != nil {
		logger.LogMessages <- logger.LogMessage{
			Message: "failed to create new chat",
			Level:   slog.LevelError,
			Attrs: []slog.Attr{
				slog.String("err", err.Error()),
				slog.Uint64("user_id", uint64(userId)),
			},
		}
		return err
	}

	logger.LogMessages <- logger.LogMessage{
		Message: "new chat has been created",
		Level:   slog.LevelInfo,
		Attrs: []slog.Attr{
			slog.Any("chat", chat),
			slog.Uint64("user_id", uint64(userId)),
		},
	}
	return c.SendStatus(fiber.StatusOK)
}
