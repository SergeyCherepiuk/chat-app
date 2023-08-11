package validation_test

import (
	"errors"
	"testing"

	"github.com/SergeyCherepiuk/chat-app/pkg/http/validation"
	"github.com/SergeyCherepiuk/chat-app/utils"
)

func Test_Invalid_CreateDirectMessageBody(t *testing.T) {
	body := validation.CreateMessageBody{Message: "   "}

	actualErr := body.Validate()
	expectedErr := errors.New("message is empty")
	expectedBody := validation.CreateMessageBody{Message: ""}

	if !utils.AreErrorsEqual(actualErr, expectedErr) {
		t.Errorf("expected: %v, actual: %v\n", expectedErr, actualErr)
	}
	if body != expectedBody {
		t.Errorf("expected: %v, actual: %v\n", expectedBody, body)
	}
}

func Test_Valid_CreateDirectMessageBody(t *testing.T) {
	body := validation.CreateMessageBody{Message: "Valid"}

	actualErr := body.Validate()
	var expectedErr error = nil
	expectedBody := validation.CreateMessageBody{Message: "Valid"}

	if !utils.AreErrorsEqual(actualErr, expectedErr) {
		t.Errorf("expected: %v, actual: %v\n", expectedErr, actualErr)
	}
	if body != expectedBody {
		t.Errorf("expected: %v, actual: %v\n", expectedBody, body)
	}
}

func Test_Invalid_UpdateDirectMessageRequestBody(t *testing.T) {
	body := validation.UpdateMessageRequestBody{Message: "   "}

	actualErr := body.Validate()
	expectedErr := errors.New("message is empty")
	expectedBody := validation.UpdateMessageRequestBody{Message: ""}

	if !utils.AreErrorsEqual(actualErr, expectedErr) {
		t.Errorf("expected: %v, actual: %v\n", expectedErr, actualErr)
	}
	if body != expectedBody {
		t.Errorf("expected: %v, actual: %v\n", expectedBody, body)
	}
}

func Test_Valid_UpdateDirectMessageRequestBody(t *testing.T) {
	body := validation.UpdateMessageRequestBody{Message: "Valid"}

	actualErr := body.Validate()
	var expectedErr error = nil
	expectedBody := validation.UpdateMessageRequestBody{Message: "Valid"}

	if !utils.AreErrorsEqual(actualErr, expectedErr) {
		t.Errorf("expected: %v, actual: %v\n", expectedErr, actualErr)
	}
	if body != expectedBody {
		t.Errorf("expected: %v, actual: %v\n", expectedBody, body)
	}
}
