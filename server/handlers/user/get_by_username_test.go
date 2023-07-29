package userhandler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func TestUnauthorizedGetByUsernameRequest(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/markwatson", nil)

	response, _ := app.Test(request)
	if response.StatusCode != fiber.StatusUnauthorized {
		t.Errorf(
			"expected status code: %v, actual status code: %v\n",
			fiber.StatusUnauthorized,
			response.StatusCode,
		)
	}
}

func TestUsernameNotFoundGetByUsernameRequest(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/theodavis", nil)
	request.AddCookie(&http.Cookie{
		Name:     "session_id",
		Value:    uuid.NewString(),
		HttpOnly: true,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
	})

	response, _ := app.Test(request)
	if response.StatusCode != fiber.StatusInternalServerError {
		t.Errorf(
			"expected status code: %v, actual status code: %v\n",
			fiber.StatusInternalServerError,
			response.StatusCode,
		)
	}
}

func TestValidGetByUsernameRequest(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/markwatson", nil)
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
			fiber.StatusUnauthorized,
			response.StatusCode,
		)
	}
}
