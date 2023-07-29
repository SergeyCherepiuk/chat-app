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

func TestValidSignUpRequestBody(t *testing.T) {
	body := authhandler.SignUpRequestBody{
		FirstName: "John",
		LastName:  "White",
		Username:  "johnwhite",
		Password:  "secret12",
	}

	actual := body.Validate()
	var expected error = nil

	if !utils.AreErrorsEqual(actual, expected) {
		t.Errorf("expected: %v, actual: %v\n", expected, actual)
	}
}

func TestEmptySignUpRequestBody(t *testing.T) {
	body := authhandler.SignUpRequestBody{}

	actual := body.Validate()
	expected := errors.Join(
		errors.New("first name is empty"),
		errors.New("last name is empty"),
		errors.New("username is empty"),
		errors.New("password is empty"),
	)

	if !utils.AreErrorsEqual(actual, expected) {
		t.Errorf("expected: %v, actual: %v\n", expected, actual)
	}
}

func TestShortPasswordSignUpRequestBody(t *testing.T) {
	body := authhandler.SignUpRequestBody{
		FirstName: "John",
		LastName:  "White",
		Username:  "johnwhite",
		Password:  "secret",
	}

	actual := body.Validate()
	expected := errors.New("password is too short")

	if !utils.AreErrorsEqual(actual, expected) {
		t.Errorf("expected: %v, actual: %v\n", expected, actual)
	}
}

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
}
