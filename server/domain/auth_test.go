package domain

import (
	"errors"
	"testing"

	"github.com/SergeyCherepiuk/chat-app/utils"
)

func TestValidSignUpRequestBody(t *testing.T) {
	body := SignUpRequestBody{
		FirstName: "John",
		LastName:  "White",
		Username:  "johnwhite",
		Password:  "secret12",
	}

	actual := body.Validate()
	var expected error = nil

	if !utils.AreErrorsEqual(actual, expected) {
		t.Errorf("Expected: %v, got: %v\n", expected, actual)
	}
}

func TestEmptySignUpRequestBody(t *testing.T) {
	body := SignUpRequestBody{}

	actual := body.Validate()
	expected := errors.Join(
		errors.New("first name is empty"),
		errors.New("last name is empty"),
		errors.New("username is empty"),
		errors.New("password is empty"),
	)

	if !utils.AreErrorsEqual(actual, expected) {
		t.Errorf("Expected: %v, got: %v\n", expected, actual)
	}
}

func TestShortPasswordSignUpRequestBody(t *testing.T) {
	body := SignUpRequestBody{
		FirstName: "John",
		LastName:  "White",
		Username:  "johnwhite",
		Password:  "secret",
	}

	actual := body.Validate()
	expected := errors.New("password is too short")

	if !utils.AreErrorsEqual(actual, expected) {
		t.Errorf("Expected: %v, got: %v\n", expected, actual)
	}
}

func TestValidLogInRequestBody(t *testing.T) {
	body := LoginRequestBody{
		Username: "johnwhite",
		Password: "secret12",
	}

	actual := body.Validate()
	var expected error = nil

	if !utils.AreErrorsEqual(actual, expected) {
		t.Errorf("Expected: %v, got: %v\n", expected, actual)
	}
}

func TestEmptyLoginRequestBody(t *testing.T) {
	body := LoginRequestBody{}

	actual := body.Validate()
	expected := errors.Join(
		errors.New("username is empty"),
		errors.New("password is empty"),
	)

	if !utils.AreErrorsEqual(actual, expected) {
		t.Errorf("Expected: %v, got: %v\n", expected, actual)
	}
}

func TestShortPasswordLoginRequestBody(t *testing.T) {
	body := LoginRequestBody{
		Username: "johnwhite",
		Password: "secret",
	}

	actual := body.Validate()
	expected := errors.New("password is too short")

	if !utils.AreErrorsEqual(actual, expected) {
		t.Errorf("Expected: %v, got: %v\n", expected, actual)
	}
}
