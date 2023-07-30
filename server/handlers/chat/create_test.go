package chathandler_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	chathandler "github.com/SergeyCherepiuk/chat-app/handlers/chat"
	"github.com/SergeyCherepiuk/chat-app/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func TestValidCreateChatRequestBody(t *testing.T) {
	body := chathandler.CreateChatRequestBody{
		Name: "New chat",
	}

	actual := body.Validate()
	var expected error = nil

	if !utils.AreErrorsEqual(actual, expected) {
		t.Errorf("expected: %v, got: %v", actual, expected)
	}
}

func TestEmptyCreateChatRequestBody(t *testing.T) {
	body := chathandler.CreateChatRequestBody{}

	actual := body.Validate()
	expected := errors.New("name is empty")

	if !utils.AreErrorsEqual(actual, expected) {
		t.Errorf("expected: %v, got: %v", actual, expected)
	}
}

func TestUnauthorizedCreateChatRequest(t *testing.T) {
	request := httptest.NewRequest(http.MethodPost, "/", nil)

	response, _ := app.Test(request)
	if response.StatusCode != fiber.StatusUnauthorized {
		t.Errorf(
			"expected status code: %v, actual status code: %v\n",
			fiber.StatusUnauthorized,
			response.StatusCode,
		)
	}
}

func TestInvalidCreateChatRequest(t *testing.T) {
	request := httptest.NewRequest(http.MethodPost, "/", nil)
	request.AddCookie(&http.Cookie{
		Name:     "session_id",
		Value:    uuid.NewString(),
		HttpOnly: true,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
	})

	response, _ := app.Test(request)
	if response.StatusCode != fiber.StatusUnprocessableEntity {
		t.Errorf(
			"expected status code: %v, actual status code: %v\n",
			fiber.StatusUnprocessableEntity,
			response.StatusCode,
		)
	}
}

func TestValidCreateChatRequest(t *testing.T) {
	request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"name": "New chat"}`))
	request.Header.Set("Content-Type", "application/json")
	request.AddCookie(&http.Cookie{
		Name:     "session_id",
		Value:    uuid.NewString(),
		HttpOnly: true,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
	})

	response, _ := app.Test(request)
	if response.StatusCode != fiber.StatusOK {
		t.Errorf(
			"expected status code: %v, actual status code: %v\n",
			fiber.StatusOK,
			response.StatusCode,
		)
	}
}
