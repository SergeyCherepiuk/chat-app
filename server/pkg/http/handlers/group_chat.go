package handlers

import (
	"strconv"

	"github.com/SergeyCherepiuk/chat-app/domain"
	"github.com/SergeyCherepiuk/chat-app/pkg/http/validation"
	"github.com/SergeyCherepiuk/chat-app/pkg/log"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/exp/slog"
)

type GroupChatHandler struct {
	groupChatService domain.GroupChatService
}

func NewGroupChatHandler(groupChatService domain.GroupChatService) *GroupChatHandler {
	return &GroupChatHandler{groupChatService: groupChatService}
}

func (handler GroupChatHandler) GetChat(c *fiber.Ctx) error {
	log := log.Logger{}

	userId := c.Locals("user_id").(uint)
	log.With(slog.Uint64("user_id", uint64(userId)))

	chatId, err := strconv.ParseUint(c.Params("chat_id"), 10, 64)
	if err != nil {
		log.Error(
			"failed to parse group chat id",
			slog.String("err", err.Error()),
			slog.String("chat_id", c.Params("chat_id")),
		)
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	log.With(slog.Uint64("chat_id", uint64(chatId)))

	chat, err := handler.groupChatService.GetChat(uint(chatId))
	if err != nil {
		slog.Error("failed to get the group chat info", slog.String("err", err.Error()))
		return err
	}

	responseBody := validation.GetGroupChatResponseBody{
		Name: chat.Name,
	}
	log.Info("group chat info has been sent to the user")
	return c.JSON(responseBody)
}

func (handler GroupChatHandler) Create(c *fiber.Ctx) error {
	log := log.Logger{}

	userId := c.Locals("user_id").(uint)
	log.With(slog.Uint64("user_id", uint64(userId)))

	body := validation.CreateGroupChatRequestBody{}
	if err := c.BodyParser(&body); err != nil {
		log.Error("failed to parse request body", slog.String("err", err.Error()))
		return err
	}

	if err := body.Validate(); err != nil {
		log.Error("invalid request body", slog.String("err", err.Error()))
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	chat := domain.GroupChat{
		Name:      body.Name,
		CreatorID: userId,
	}
	if err := handler.groupChatService.Create(&chat); err != nil {
		log.Error("failed to store the group chat", slog.String("err", err.Error()))
	}

	log.Info("group chat has been stored", slog.Any("chat", chat))
	return nil
}

func (handler GroupChatHandler) Update(c *fiber.Ctx) error {
	log := log.Logger{}

	userId := c.Locals("user_id").(uint)
	log.With(slog.Uint64("user_id", uint64(userId)))

	chatId := c.Locals("chat_id").(uint)
	log.With(slog.Uint64("chat_id", uint64(chatId)))

	body := validation.UpdateGroupChatRequestBody{}
	if err := c.BodyParser(&body); err != nil {
		log.Error("failed to parse request body", slog.String("err", err.Error()))
		return err
	}

	updates := body.ToMap()
	log.With(slog.Any("updates", updates))

	if err := handler.groupChatService.Update(chatId, updates); err != nil {
		log.Error("failed to update the group chat", slog.String("err", err.Error()))
		return err
	}

	slog.Info("group chat has been updated")
	return c.SendStatus(fiber.StatusOK)
}

func (handler GroupChatHandler) Delete(c *fiber.Ctx) error {
	log := log.Logger{}

	userId := c.Locals("user_id").(uint)
	log.With(slog.Uint64("user_id", uint64(userId)))

	chatId := c.Locals("chat_id").(uint)
	log.With(slog.Uint64("chat_id", uint64(chatId)))

	if err := handler.groupChatService.Delete(chatId); err != nil {
		log.Error("failed to delete the group chat", slog.String("err", err.Error()))
		return err
	}

	log.Info("group chat has been deleted")
	return c.SendStatus(fiber.StatusOK)
}
