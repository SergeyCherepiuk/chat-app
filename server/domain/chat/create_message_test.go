package chatdomain_test

import (
	"errors"
	"testing"

	chatdomain "github.com/SergeyCherepiuk/chat-app/domain/chat"
	"github.com/SergeyCherepiuk/chat-app/utils"
)

func TestInvalidCreateMessageBody(t *testing.T) {
	body := chatdomain.CreateMessageBody{Message: "   "}

	actualErr := body.Validate()
	expectedErr := errors.New("message is empty")
	expectedBody := chatdomain.CreateMessageBody{Message: ""}

	if !utils.AreErrorsEqual(actualErr, expectedErr) {
		t.Errorf("expected: %v, actual: %v\n", expectedErr, actualErr)
	}
	if body != expectedBody {
		t.Errorf("expected: %v, actual: %v\n", expectedBody, body)
	}
}

func TestValidCreateMessageBody(t *testing.T) {
	body := chatdomain.CreateMessageBody{Message: "Valid"}

	actualErr := body.Validate()
	var expectedErr error = nil
	expectedBody := chatdomain.CreateMessageBody{Message: "Valid"}

	if !utils.AreErrorsEqual(actualErr, expectedErr) {
		t.Errorf("expected: %v, actual: %v\n", expectedErr, actualErr)
	}
	if body != expectedBody {
		t.Errorf("expected: %v, actual: %v\n", expectedBody, body)
	}
}