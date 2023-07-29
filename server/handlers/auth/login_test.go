package authhandler_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	authhandler "github.com/SergeyCherepiuk/chat-app/handlers/auth"
	"github.com/SergeyCherepiuk/chat-app/utils"
	"github.com/gofiber/fiber/v2"
)

func TestValidLoginRequestBody(t *testing.T) {
	body := authhandler.LoginRequestBody{
		Username: "johnwhite",
		Password: "secret12",
	}

	actual := body.Validate()
	var expected error = nil

	if !utils.AreErrorsEqual(actual, expected) {
		t.Errorf("expected: %v, actual: %v\n", expected, actual)
	}
}

func TestEmptyLoginRequestBody(t *testing.T) {
	body := authhandler.LoginRequestBody{}

	actual := body.Validate()
	expected := errors.Join(
		errors.New("username is empty"),
		errors.New("password is empty"),
	)

	if !utils.AreErrorsEqual(actual, expected) {
		t.Errorf("expected: %v, actual: %v\n", expected, actual)
	}
}

func TestShortPasswordLoginRequestBody(t *testing.T) {
	body := authhandler.LoginRequestBody{
		Username: "johnwhite",
		Password: "secret",
	}

	actual := body.Validate()
	expected := errors.New("password is too short")

	if !utils.AreErrorsEqual(actual, expected) {
		t.Errorf("expected: %v, actual: %v\n", expected, actual)
	}
}

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

func TestInvalidLoginRequest(t *testing.T) {
	request := httptest.NewRequest(
		http.MethodPost,
		"/login",
		strings.NewReader(`{
			"username": "johndoe",
			"password": "sec"
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

func TestValidLoginRequest(t *testing.T) {
	request := httptest.NewRequest(
		http.MethodPost,
		"/login",
		strings.NewReader(`{
			"username": "johndoe",
			"password": "secret123"
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
