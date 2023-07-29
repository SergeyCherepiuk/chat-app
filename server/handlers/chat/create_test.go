package chathandler_test

import (
	"errors"
	"testing"

	chathandler "github.com/SergeyCherepiuk/chat-app/handlers/chat"
	"github.com/SergeyCherepiuk/chat-app/utils"
)

func TestValidCreateChatRequestBody(t *testing.T) {
	body := chathandler.CreateChatRequestBody{
		Name: "New chat",
	}

	actual := body.Validate()
	var expected error = nil

	if !utils.AreErrorsEqual(actual, expected) {
		t.Errorf("expected: %v, got: %v", actual, expected)
	}
}

func TestEmptyCreateChatRequestBody(t *testing.T) {
	body := chathandler.CreateChatRequestBody{}

	actual := body.Validate()
	expected := errors.New("name is empty")

	if !utils.AreErrorsEqual(actual, expected) {
		t.Errorf("expected: %v, got: %v", actual, expected)
	}
}
