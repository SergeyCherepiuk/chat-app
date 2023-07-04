package handlers

import (
	"asdkoda/session-auth/src/models"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthHandler struct {
	db *gorm.DB
}

func NewAuthHandler(db *gorm.DB) AuthHandler {
	return AuthHandler{db: db}
}

func createCookie(c *fiber.Ctx, sessionId string, expiresIn time.Duration) {
	cookie := new(fiber.Cookie)
	cookie.Name = "session_id"
	cookie.Value = sessionId
	cookie.Expires = time.Now().Add(expiresIn)
	cookie.HTTPOnly = true
	c.Cookie(cookie)
}

type authRequestBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (body authRequestBody) validate() error {
	if strings.TrimSpace(body.Username) == "" {
		return errors.New("invalid username")
	}

	if body.Password != strings.TrimSpace(body.Password) {
		return errors.New("password has whitespaces")
	}

	if len(body.Password) < 8 {
		return errors.New("password is too short (minimum 8 characters required)")
	}

	return nil
}

func (handler AuthHandler) SingUp(c *fiber.Ctx) error {
	body := authRequestBody{}
	if err := c.BodyParser(&body); err != nil {
		return err
	}

	err := body.validate()
	if err != nil {
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		return err
	}
	user := models.User{Username: body.Username, Password: string(hashedPassword)}

	if result := handler.db.Create(&user); result.Error != nil {
		return result.Error
	}

	session := models.Session{UserID: user.ID, ExpiresAt: time.Now().Add(7 * 24 * time.Hour)}
	if result := handler.db.Create(&session); result.Error != nil {
		return result.Error
	}

	createCookie(c, fmt.Sprint(session.ID), 7*24*time.Hour)

	return c.SendString("Signed up successfully")
}
