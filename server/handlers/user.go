package handlers

import (
	"time"

	"github.com/SergeyCherepiuk/chat-app/domain"
	"github.com/SergeyCherepiuk/chat-app/storage"
	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	storage *storage.UserStorage
}

func NewUserHandler(storage *storage.UserStorage) *UserHandler {
	return &UserHandler{storage: storage}
}

func (handler UserHandler) GetMe(c *fiber.Ctx) error {
	userId, ok := c.Locals("user_id").(uint)
	if !ok {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	user, err := handler.storage.GetById(userId)
	if err != nil {
		return err
	}

	responseBody := domain.GetUserResponseBody{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Username:  user.Username,
	}
	return c.JSON(responseBody)
}

func (handler UserHandler) GetUser(c *fiber.Ctx) error {
	user, err := handler.storage.GetByUsername(c.Params("username"))
	if err != nil {
		return err
	}

	responseBody := domain.GetUserResponseBody{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Username:  user.Username,
	}
	return c.JSON(responseBody)
}

func (handler UserHandler) UpdateMe(c *fiber.Ctx) error {
	userId, ok := c.Locals("user_id").(uint)
	if !ok {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	body := domain.UpdateUserRequestBody{}
	if err := c.BodyParser(&body); err != nil {
		return err
	}

	updates := body.ToMap()
	if err := handler.storage.Update(userId, updates); err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusOK)
}

func (handler UserHandler) DeleteMe(c *fiber.Ctx) error {
	userId, ok := c.Locals("user_id").(uint)
	if !ok {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	if err := handler.storage.Delete(userId); err != nil {
		return err
	}

	c.Cookie(&fiber.Cookie{
		Name:    "session_id",
		Expires: time.Now(),
	})
	return c.SendStatus(fiber.StatusOK)
}
