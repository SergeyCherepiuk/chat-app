package chathandler_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func Test_Unauthorized_UpdateMessageRequest(t *testing.T) {
	request := httptest.NewRequest(http.MethodPut, "/chat/markwatson/1", nil)

	response, _ := app.Test(request)
	if response.StatusCode != fiber.StatusUnauthorized {
		t.Errorf(
			"expected status code: %v, actual status code: %v\n",
			fiber.StatusUnauthorized,
			response.StatusCode,
		)
	}
}

func Test_UnparsableBody_UpdateMessageRequest(t *testing.T) {
	request := httptest.NewRequest(http.MethodPut, "/chat/markwatson/1", nil)
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

func Test_InvalidBody_UpdateMessageRequest(t *testing.T) {
	request := httptest.NewRequest(
		http.MethodPut,
		"/chat/markwatson/1",
		strings.NewReader(`{
			"message": ""
		}`),
	)
	request.Header.Set("Content-Type", "application/json")
	request.AddCookie(&http.Cookie{
		Name:     "session_id",
		Value:    uuid.NewString(),
		HttpOnly: true,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
	})

	response, _ := app.Test(request)
	if response.StatusCode != fiber.StatusBadRequest {
		t.Errorf(
			"expected status code: %v, actual status code: %v\n",
			fiber.StatusBadRequest,
			response.StatusCode,
		)
	}
}

func Test_UserNotAnAuthor_UpdateMessageRequest(t *testing.T) {
	request := httptest.NewRequest(
		http.MethodPut,
		"/chat/markwatson/2",
		strings.NewReader(`{
			"message": "New message"
		}`),
	)
	request.Header.Set("Content-Type", "application/json")
	request.AddCookie(&http.Cookie{
		Name:     "session_id",
		Value:    uuid.NewString(),
		HttpOnly: true,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
	})

	response, _ := app.Test(request)
	if response.StatusCode != fiber.StatusUnauthorized {
		t.Errorf(
			"expected status code: %v, actual status code: %v\n",
			fiber.StatusUnauthorized,
			response.StatusCode,
		)
	}
}

func Test_Valid_UpdateMessageRequest(t *testing.T) {
	request := httptest.NewRequest(
		http.MethodPut,
		"/chat/markwatson/1",
		strings.NewReader(`{
			"message": "New message"
		}`),
	)
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
