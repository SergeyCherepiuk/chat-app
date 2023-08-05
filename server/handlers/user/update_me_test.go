package userhandler_test

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	userdomain "github.com/SergeyCherepiuk/chat-app/domain/user"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func TestConvertValidUpdateUserRequestBodyToMap(t *testing.T) {
	body := userdomain.UpdateUserRequestBody{
		FirstName: "Andrew",
		LastName:  "Brown",
		Username:  "andrewbrown",
	}

	actual := body.ToMap()
	expected := map[string]any{
		"first_name": "Andrew",
		"last_name":  "Brown",
		"username":   "andrewbrown",
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected: %v, got: %v", actual, expected)
	}
}

func TestConvertEmptyUpdateUserRequestBodyToMap(t *testing.T) {
	body := userdomain.UpdateUserRequestBody{}

	actual := body.ToMap()
	expected := map[string]any{}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected: %v, got: %v", actual, expected)
	}
}

func TestConvertWhiteSpaceUpdateUserRequestBodyToMap(t *testing.T) {
	body := userdomain.UpdateUserRequestBody{
		FirstName: "",
		LastName:  " ",
		Username:  "  ",
	}

	actual := body.ToMap()
	expected := map[string]any{}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected: %v, got: %v", actual, expected)
	}
}

func TestUnauthorizedUpdateMeRequest(t *testing.T) {
	request := httptest.NewRequest(http.MethodPut, "/me", nil)

	response, _ := app.Test(request)
	if response.StatusCode != fiber.StatusUnauthorized {
		t.Errorf(
			"expected status code: %v, actual status code: %v\n",
			fiber.StatusUnauthorized,
			response.StatusCode,
		)
	}
}

func TestInvalidUpdateMeRequest(t *testing.T) {
	request := httptest.NewRequest(http.MethodPut, "/me", nil)
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

func TestValidUpdateMeRequest(t *testing.T) {
	request := httptest.NewRequest(
		http.MethodPut,
		"/me",
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
