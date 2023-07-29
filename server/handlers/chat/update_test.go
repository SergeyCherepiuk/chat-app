package chathandler_test

import (
	"reflect"
	"testing"

	chathandler "github.com/SergeyCherepiuk/chat-app/handlers/chat"
)

func TestConvertValidUpdateChatRequestBodyToMap(t *testing.T) {
	body := chathandler.UpdateChatRequestBody{
		Name: "New chat's name",
	}

	actual := body.ToMap()
	expected := map[string]any{
		"name": "New chat's name",
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected: %v, got: %v", actual, expected)
	}
}

func TestConvertEmptyUpdateChatRequestBodyToMap(t *testing.T) {
	body := chathandler.UpdateChatRequestBody{}

	actual := body.ToMap()
	expected := map[string]any{}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected: %v, got: %v", actual, expected)
	}
}

func TestConvertWhiteSpaceUpdateChatRequestBodyToMap(t *testing.T) {
	body := chathandler.UpdateChatRequestBody{
		Name: "   ",
	}

	actual := body.ToMap()
	expected := map[string]any{}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected: %v, got: %v", actual, expected)
	}
}
