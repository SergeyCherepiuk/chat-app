package handlers

import (
	"time"

	"github.com/SergeyCherepiuk/chat-app/domain"
	"github.com/SergeyCherepiuk/chat-app/pkg/http/validation"
	"github.com/SergeyCherepiuk/chat-app/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/exp/slog"
)

type UserHandler struct {
	userService domain.UserService
}

func NewUserHandler(userService domain.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (handler UserHandler) GetMe(c *fiber.Ctx) error {
	log := logger.Logger{}

	userId, _ := c.Locals("user_id").(uint)
	log.With(slog.Uint64("user_id", uint64(userId)))

	user, err := handler.userService.GetById(userId)
	if err != nil {
		log.Error("failed to get user by id", slog.String("err", err.Error()))
		return err
	}

	responseBody := validation.GetUserResponseBody{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Username:  user.Username,
	}
	log.Info("user's info has been sent to the user", slog.Any("user", responseBody))
	return c.JSON(responseBody)
}

func (handler UserHandler) GetByUsername(c *fiber.Ctx) error {
	log := logger.Logger{}

	userId, _ := c.Locals("user_id").(uint)
	log.With(slog.Uint64("user_id", uint64(userId)))

	user, err := handler.userService.GetByUsername(c.Params("username"))
	if err != nil {
		log.Error(
			"failed to get user by username",
			slog.String("err", err.Error()),
			slog.String("username", c.Params("username")),
		)
		return err
	}

	responseBody := validation.GetUserResponseBody{
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Username:    user.Username,
		Description: user.Description,
	}
	log.Info("user's info has been sent to the user", slog.Any("user", responseBody))
	return c.JSON(responseBody)
}

func (handler UserHandler) UpdateMe(c *fiber.Ctx) error {
	log := logger.Logger{}

	userId, _ := c.Locals("user_id").(uint)
	log.With(slog.Uint64("user_id", uint64(userId)))

	body := validation.UpdateUserRequestBody{}
	if err := c.BodyParser(&body); err != nil {
		log.Error("failed to parse request body", slog.String("err", err.Error()))
		return err
	}

	updates := body.ToMap()
	if err := handler.userService.Update(userId, updates); err != nil {
		log.Error(
			"failed to update the user",
			slog.String("err", err.Error()),
			slog.Any("updates", updates),
		)
		return err
	}

	log.Info("user has been updated", slog.Any("updates", updates))
	return c.SendStatus(fiber.StatusOK)
}

func (handler UserHandler) DeleteMe(c *fiber.Ctx) error {
	log := logger.Logger{}

	userId, _ := c.Locals("user_id").(uint)
	log.With(slog.Uint64("user_id", uint64(userId)))

	if err := handler.userService.Delete(userId); err != nil {
		log.Error("failed to delete the user", slog.String("err", err.Error()))
		return err
	}

	c.Cookie(&fiber.Cookie{
		Name:    "session_id",
		Expires: time.Now(),
	})
	log.Info("user has been deleted")
	return c.SendStatus(fiber.StatusOK)
}
