package authhandler_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestUnparsableSignUpRequest(t *testing.T) {
	request := httptest.NewRequest(http.MethodPost, "/signup", nil)

	response, _ := app.Test(request)
	if response.StatusCode != fiber.StatusUnprocessableEntity {
		t.Errorf(
			"expected status code: %v, actual status code: %v\n",
			fiber.StatusUnprocessableEntity,
			response.StatusCode,
		)
	}
}

func TestInvalidSignUpRequest(t *testing.T) {
	request := httptest.NewRequest(
		http.MethodPost,
		"/signup",
		strings.NewReader(`{
			"first_name": "John"
		}`),
	)
	request.Header.Set("Content-Type", "application/json")

	response, _ := app.Test(request)
	if response.StatusCode != fiber.StatusBadRequest {
		t.Errorf(
			"expected status code: %v, actual status code: %v\n",
			fiber.StatusBadRequest,
			response.StatusCode,
		)
	}
}

func TestValidSignUpRequest(t *testing.T) {
	request := httptest.NewRequest(
		http.MethodPost,
		"/signup",
		strings.NewReader(`{
			"first_name": "John",
			"last_name": "Doe",
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
}
