package authhandler_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestUnparsableLoginRequest(t *testing.T) {
	request := httptest.NewRequest(http.MethodPost, "/login", nil)

	response, _ := app.Test(request)
	if response.StatusCode != fiber.StatusUnprocessableEntity {
		t.Errorf(
			"expected status code: %v, actual status code: %v\n",
			fiber.StatusUnprocessableEntity,
			response.StatusCode,
		)
	}
}

func TestValidLoginRequest(t *testing.T) {
	request := httptest.NewRequest(
		http.MethodPost,
		"/login",
		strings.NewReader(`{
			"username": "johndoe",
			"password": "Secret123!"
		}`),
	)
	request.Header.Set("Content-Type", "application/json")

	response, _ := app.Test(request)
	if response.StatusCode != fiber.StatusOK {
		t.Errorf(
			"expected status code: %v, actual status code: %v\n",
			fiber.StatusOK,
			response.StatusCode,
		)
	}
	for _, c := range response.Cookies() {
		if c.Name == "session_id" && c.HttpOnly == true {
			break
		}
		t.Errorf("cookie hasn't been set")
	}
}
