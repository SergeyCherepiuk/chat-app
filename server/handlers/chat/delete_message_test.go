package chathandler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func Test_Unauthorized_DeleteMessageRequest(t *testing.T) {
	request := httptest.NewRequest(http.MethodDelete, "/chat/johndoe/1", nil)

	response, _ := app.Test(request)
	if response.StatusCode != fiber.StatusUnauthorized {
		t.Errorf(
			"expected status code: %v, actual status code: %v\n",
			fiber.StatusUnauthorized,
			response.StatusCode,
		)
	}
}

func Test_UserNotFound_DeleteMessageRequest(t *testing.T) {
	request := httptest.NewRequest(http.MethodDelete, "/chat/unknownusername/1", nil)
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

func Test_MessageNotBelongsToChat_DeleteMessageRequest(t *testing.T) {
	request := httptest.NewRequest(http.MethodDelete, "/chat/markwatson/3", nil)
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

func Test_Valid_DeleteMessageRequest(t *testing.T) {
	request := httptest.NewRequest(http.MethodDelete, "/chat/markwatson/2", nil)
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
