package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func Test_Unparsable_SignUpRequest(t *testing.T) {
	request := httptest.NewRequest(http.MethodPost, "/api/auth/signup", nil)

	response, _ := app.Test(request)
	if response.StatusCode != fiber.StatusUnprocessableEntity {
		t.Errorf(
			"expected status code: %v, actual status code: %v\n",
			fiber.StatusUnprocessableEntity,
			response.StatusCode,
		)
	}
}

func Test_Invalid_SignUpRequest(t *testing.T) {
	request := httptest.NewRequest(
		http.MethodPost,
		"/api/auth/signup",
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

func Test_Valid_SignUpRequest(t *testing.T) {
	request := httptest.NewRequest(
		http.MethodPost,
		"/api/auth/signup",
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

func Test_Unparsable_LoginRequest(t *testing.T) {
	request := httptest.NewRequest(http.MethodPost, "/api/auth/login", nil)

	response, _ := app.Test(request)
	if response.StatusCode != fiber.StatusUnprocessableEntity {
		t.Errorf(
			"expected status code: %v, actual status code: %v\n",
			fiber.StatusUnprocessableEntity,
			response.StatusCode,
		)
	}
}

func Test_Valid_LoginRequest(t *testing.T) {
	request := httptest.NewRequest(
		http.MethodPost,
		"/api/auth/login",
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
		if c.Name == "session_id" && c.HttpOnly {
			break
		}
		t.Errorf("cookie hasn't been set")
	}
}

func Test_Invalid_LogoutRequest(t *testing.T) {
	request := httptest.NewRequest(http.MethodPost, "/api/auth/logout", nil)

	response, _ := app.Test(request)
	if response.StatusCode != fiber.StatusBadRequest {
		t.Errorf(
			"expected status code: %v, actual status code: %v\n",
			fiber.StatusBadRequest,
			response.StatusCode,
		)
	}
}

func Test_Valid_LogoutRequest(t *testing.T) {
	request := httptest.NewRequest(http.MethodPost, "/api/auth/logout", nil)
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
