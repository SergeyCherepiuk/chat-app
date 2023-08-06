package chatdomain_test

import (
	"errors"
	"testing"

	chatdomain "github.com/SergeyCherepiuk/chat-app/domain/chat"
	"github.com/SergeyCherepiuk/chat-app/utils"
)

func TestInvalidUpdateMessageRequestBody(t *testing.T) {
	body := chatdomain.UpdateMessageRequestBody{Message: "   "}

	actualErr := body.Validate()
	expectedErr := errors.New("message is empty")
	expectedBody := chatdomain.UpdateMessageRequestBody{Message: ""}

	if !utils.AreErrorsEqual(actualErr, expectedErr) {
		t.Errorf("expected: %v, actual: %v\n", expectedErr, actualErr)
	}
	if body != expectedBody {
		t.Errorf("expected: %v, actual: %v\n", expectedBody, body)
	}
}

func TestValidUpdateMessageRequestBody(t *testing.T) {
	body := chatdomain.UpdateMessageRequestBody{Message: "Valid"}

	actualErr := body.Validate()
	var expectedErr error = nil
	expectedBody := chatdomain.UpdateMessageRequestBody{Message: "Valid"}

	if !utils.AreErrorsEqual(actualErr, expectedErr) {
		t.Errorf("expected: %v, actual: %v\n", expectedErr, actualErr)
	}
	if body != expectedBody {
		t.Errorf("expected: %v, actual: %v\n", expectedBody, body)
	}
}