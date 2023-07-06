package handlers

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/SergeyCherepiuk/session-auth/src/models"
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

func createCookie(c *fiber.Ctx, sessionId string) {
	cookie := new(fiber.Cookie)
	cookie.Name = "session_id"
	cookie.Value = sessionId
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

	session := models.Session{UserID: user.ID, ExpiresAt: time.Now().Add(10 * time.Second)}
	if result := handler.db.Create(&session); result.Error != nil {
		return result.Error
	}

	createCookie(c, fmt.Sprint(session.ID))
	return c.SendString("Signed up successfully")
}

func (handler AuthHandler) Login(c *fiber.Ctx) error {
	body := authRequestBody{}
	if err := c.BodyParser(&body); err != nil {
		return err
	}

	err := body.validate()
	if err != nil {
		return err
	}

	user := models.User{}
	if r := handler.db.Where("username = ?", body.Username).First(&user); r.Error != nil {
		return r.Error
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)) != nil {
		return errors.New("wrong password")
	}

	if r := handler.db.Where("user_id = ?", user.ID).Delete(&models.Session{}); r.Error != nil {
		return r.Error
	}

	session := models.Session{UserID: user.ID, ExpiresAt: time.Now().Add(10 * time.Second)}
	if result := handler.db.Create(&session); result.Error != nil {
		return result.Error
	}

	createCookie(c, fmt.Sprint(session.ID))
	return c.SendString("Logged in successfully")
}

func (handler AuthHandler) Logout(c *fiber.Ctx) error {
	sessionId, err := strconv.ParseUint(c.Cookies("session_id", ""), 10, 64)
	if err != nil {
		return err
	}

	session := models.Session{}
	if r := handler.db.First(&session, sessionId); r.Error != nil || r.RowsAffected < 1 {
		return errors.New("session not found")
	}

	if r := handler.db.Where("user_id = ?", session.UserID).Delete(&models.Session{}); r.Error != nil {
		return r.Error
	}

	return c.Status(fiber.StatusOK).SendString("Logged out successfully")
}
