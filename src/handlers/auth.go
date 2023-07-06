package handlers

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/SergeyCherepiuk/session-auth/src/auth"
	"github.com/SergeyCherepiuk/session-auth/src/models"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthHandler struct {
	pdb            *gorm.DB
	rdb            *redis.Client
	sessionManager *auth.SessionManager
}

func NewAuthHandler(pdb *gorm.DB, rdb *redis.Client, sessionManager *auth.SessionManager) *AuthHandler {
	return &AuthHandler{pdb: pdb, rdb: rdb, sessionManager: sessionManager}
}

func createCookie(c *fiber.Ctx, sessionId string) {
	cookie := fiber.Cookie{
		Name:     "session_id",
		Value:    sessionId,
		HTTPOnly: true,
	}
	c.Cookie(&cookie)
}

func deleteCookie(c *fiber.Ctx) {
	cookie := fiber.Cookie{
		Name: "session_id",
		Expires: time.Now(),
	}
	c.Cookie(&cookie)
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

	if result := handler.pdb.Create(&user); result.Error != nil {
		return result.Error
	}

	session, err := handler.sessionManager.CreateSession(user.ID)
	if err != nil {
		return err
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
	if r := handler.pdb.Where("username = ?", body.Username).First(&user); r.Error != nil {
		return r.Error
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)) != nil {
		return errors.New("wrong password")
	}

	session, err := handler.sessionManager.CreateSession(user.ID)
	if err != nil {
		return err
	}

	createCookie(c, fmt.Sprint(session.ID))
	return c.SendString("Logged in successfully")
}

func (handler AuthHandler) Logout(c *fiber.Ctx) error {
	// TODO: Delete session with sessionManager (by sessionId from the cookies)

	deleteCookie(c)
	return c.SendString("Logged out successfully")
}
