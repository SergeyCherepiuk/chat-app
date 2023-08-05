package authdomain_test

import (
	"errors"
	"testing"

	authdomain "github.com/SergeyCherepiuk/chat-app/domain/auth"
	"github.com/SergeyCherepiuk/chat-app/utils"
)

func TestValidSignUpRequestBody(t *testing.T) {
	body := authdomain.SignUpRequestBody{
		FirstName: "John",
		LastName:  "White",
		Username:  "johnwhite",
		Password:  "Secret12!",
	}

	actual := body.Validate()
	var expected error = nil

	if !utils.AreErrorsEqual(actual, expected) {
		t.Errorf("expected: %v, actual: %v\n", expected, actual)
	}
}

func TestEmptySignUpRequestBody(t *testing.T) {
	body := authdomain.SignUpRequestBody{}

	actual := body.Validate()
	expected := errors.Join(
		errors.New("invalid first name"),
		errors.New("invalid last name"),
		errors.New("invalid username"),
		errors.New("password must be at least 8 characters long"),
	)

	if !utils.AreErrorsEqual(actual, expected) {
		t.Errorf("expected: %v, actual: %v\n", expected, actual)
	}
}

func TestShortPasswordSignUpRequestBody(t *testing.T) {
	body := authdomain.SignUpRequestBody{
		FirstName: "John",
		LastName:  "White",
		Username:  "johnwhite",
		Password:  "secret",
	}

	actual := body.Validate()
	expected := errors.New("password must be at least 8 characters long")

	if !utils.AreErrorsEqual(actual, expected) {
		t.Errorf("expected: %v, actual: %v\n", expected, actual)
	}
}

func TestPasswordMustContainLowercaseSignUpRequestBody(t *testing.T) {
	body := authdomain.SignUpRequestBody{
		FirstName: "John",
		LastName:  "White",
		Username:  "johnwhite",
		Password:  "12345678",
	}

	actual := body.Validate()
	expected := errors.New("password must contain at least one lowercase character")

	if !utils.AreErrorsEqual(actual, expected) {
		t.Errorf("expected: %v, actual: %v\n", expected, actual)
	}
}

func TestPasswordMustContainUppercaseSignUpRequestBody(t *testing.T) {
	body := authdomain.SignUpRequestBody{
		FirstName: "John",
		LastName:  "White",
		Username:  "johnwhite",
		Password:  "secret12",
	}

	actual := body.Validate()
	expected := errors.New("password must contain at least one uppercase character")

	if !utils.AreErrorsEqual(actual, expected) {
		t.Errorf("expected: %v, actual: %v\n", expected, actual)
	}
}

func TestPasswordMustContainDigitSignUpRequestBody(t *testing.T) {
	body := authdomain.SignUpRequestBody{
		FirstName: "John",
		LastName:  "White",
		Username:  "johnwhite",
		Password:  "StrongPassword",
	}

	actual := body.Validate()
	expected := errors.New("password must contain at least one digit")

	if !utils.AreErrorsEqual(actual, expected) {
		t.Errorf("expected: %v, actual: %v\n", expected, actual)
	}
}

func TestPasswordMustContainSpecialSignUpRequestBody(t *testing.T) {
	body := authdomain.SignUpRequestBody{
		FirstName: "John",
		LastName:  "White",
		Username:  "johnwhite",
		Password:  "StrongPassword123",
	}

	actual := body.Validate()
	expected := errors.New("password must contain at least one special character")

	if !utils.AreErrorsEqual(actual, expected) {
		t.Errorf("expected: %v, actual: %v\n", expected, actual)
	}
}
