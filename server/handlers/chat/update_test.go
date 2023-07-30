package chathandler_test

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	chathandler "github.com/SergeyCherepiuk/chat-app/handlers/chat"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func TestConvertValidUpdateChatRequestBodyToMap(t *testing.T) {
	body := chathandler.UpdateChatRequestBody{
		Name: "New chat's name",
	}

	actual := body.ToMap()
	expected := map[string]any{
		"name": "New chat's name",
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected: %v, got: %v", actual, expected)
	}
}

func TestConvertEmptyUpdateChatRequestBodyToMap(t *testing.T) {
	body := chathandler.UpdateChatRequestBody{}

	actual := body.ToMap()
	expected := map[string]any{}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected: %v, got: %v", actual, expected)
	}
}

func TestConvertWhiteSpaceUpdateChatRequestBodyToMap(t *testing.T) {
	body := chathandler.UpdateChatRequestBody{
		Name: "   ",
	}

	actual := body.ToMap()
	expected := map[string]any{}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected: %v, got: %v", actual, expected)
	}
}

func TestUnauthorizedUpdateChatRequest(t *testing.T) {
	request := httptest.NewRequest(http.MethodPut, "/1", nil)

	response, _ := app.Test(request)
	if response.StatusCode != fiber.StatusUnauthorized {
		t.Errorf(
			"expected status code: %v, actual status code: %v\n",
			fiber.StatusUnauthorized,
			response.StatusCode,
		)
	}
}

func TestInvalidUpdateChatRequest(t *testing.T) {
	request := httptest.NewRequest(http.MethodPut, "/asd", nil)
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

func TestInvalidBodyUpdateChatRequest(t *testing.T) {
	request := httptest.NewRequest(http.MethodPut, "/1", nil)
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

func TestChatNotFoundBodyUpdateChatRequest(t *testing.T) {
	request := httptest.NewRequest(http.MethodPut, "/2", strings.NewReader(`{"name": "New chat's name"}`))
	request.Header.Set("Content-Type", "application/json")
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

func TestValidBodyUpdateChatRequest(t *testing.T) {
	request := httptest.NewRequest(http.MethodPut, "/1", strings.NewReader(`{"name": "New chat's name"}`))
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
