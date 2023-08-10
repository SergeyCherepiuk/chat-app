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

func Test_Unauthorized_GetMeRequest(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/api/user/me", nil)

	response, _ := app.Test(request)
	if response.StatusCode != fiber.StatusUnauthorized {
		t.Errorf(
			"expected status code: %v, actual status code: %v\n",
			fiber.StatusUnauthorized,
			response.StatusCode,
		)
	}
}

func Test_Valid_GetMeRequest(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/api/user/me", nil)
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

func Test_Unauthorized_GetByUsernameRequest(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/api/user/markwatson", nil)

	response, _ := app.Test(request)
	if response.StatusCode != fiber.StatusUnauthorized {
		t.Errorf(
			"expected status code: %v, actual status code: %v\n",
			fiber.StatusUnauthorized,
			response.StatusCode,
		)
	}
}

func Test_UserNotFound_GetByUsernameRequest(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/api/user/theodavis", nil)
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

func Test_Valid_GetByUsernameRequest(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/api/user/markwatson", nil)
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

func Test_Unauthorized_UpdateMeRequest(t *testing.T) {
	request := httptest.NewRequest(http.MethodPut, "/api/user/me", nil)

	response, _ := app.Test(request)
	if response.StatusCode != fiber.StatusUnauthorized {
		t.Errorf(
			"expected status code: %v, actual status code: %v\n",
			fiber.StatusUnauthorized,
			response.StatusCode,
		)
	}
}

func Test_Invalid_UpdateMeRequest(t *testing.T) {
	request := httptest.NewRequest(http.MethodPut, "/api/user/me", nil)
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

func Test_Valid_UpdateMeRequest(t *testing.T) {
	request := httptest.NewRequest(
		http.MethodPut,
		"/api/user/me",
		strings.NewReader(`{
			"first_name": "Jonathan",
			"last_name": "Von", 
			"username": "jonathanvon"
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
			fiber.StatusUnprocessableEntity,
			response.StatusCode,
		)
	}
}

func Test_Unauthorized_DeleteMeRequest(t *testing.T) {
	request := httptest.NewRequest(http.MethodDelete, "/api/user/me", nil)

	response, _ := app.Test(request)
	if response.StatusCode != fiber.StatusUnauthorized {
		t.Errorf(
			"expected status code: %v, actual status code: %v\n",
			fiber.StatusUnauthorized,
			response.StatusCode,
		)
	}
}

func Test_Valid_DeleteMeRequest(t *testing.T) {
	request := httptest.NewRequest(http.MethodDelete, "/api/user/me", nil)
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
			fiber.StatusUnprocessableEntity,
			response.StatusCode,
		)
	}
}
