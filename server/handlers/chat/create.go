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
	log := logger.Logger{}

	userId, _ := c.Locals("user_id").(uint)
	log.With(slog.Uint64("user_id", uint64(userId)))

	body := CreateChatRequestBody{}
	if err := c.BodyParser(&body); err != nil {
		log.Error("failed to parse request body", slog.String("err", err.Error()))
		return err
	}

	chat := models.Chat{Name: body.Name}
	if err := handler.storage.CreateChat(chat); err != nil {
		log.Error("failed to create new chat", slog.String("err", err.Error()))
		return err
	}

	log.Info("new chat has been created", slog.Any("chat", chat))
	return c.SendStatus(fiber.StatusOK)
}
