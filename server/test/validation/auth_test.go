package validation_test

import (
	"errors"
	"testing"

	"github.com/SergeyCherepiuk/chat-app/pkg/http/validation"
	"github.com/SergeyCherepiuk/chat-app/utils"
)

func Test_Valid_SignUpRequestBody(t *testing.T) {
	body := validation.SignUpRequestBody{
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

func Test_Empty_SignUpRequestBody(t *testing.T) {
	body := validation.SignUpRequestBody{}

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

func Test_ShortPassword_SignUpRequestBody(t *testing.T) {
	body := validation.SignUpRequestBody{
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

func Test_PasswordMustContainLowercase_SignUpRequestBody(t *testing.T) {
	body := validation.SignUpRequestBody{
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

func Test_PasswordMustContainUppercase_SignUpRequestBody(t *testing.T) {
	body := validation.SignUpRequestBody{
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

func Test_PasswordMustContainDigit_SignUpRequestBody(t *testing.T) {
	body := validation.SignUpRequestBody{
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

func Test_PasswordMustContainSpecial_SignUpRequestBody(t *testing.T) {
	body := validation.SignUpRequestBody{
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
