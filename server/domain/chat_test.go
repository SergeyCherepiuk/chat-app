package domain

import (
	"errors"
	"reflect"
	"testing"

	"github.com/SergeyCherepiuk/chat-app/utils"
)

func TestValidCreateChatRequestBody(t *testing.T) {
	body := CreateChatRequestBody{
		Name: "New chat",
	}

	actual := body.Validate()
	var expected error = nil

	if !utils.AreErrorsEqual(actual, expected) {
		t.Errorf("Expected: %v, got: %v", actual, expected)
	}
}

func TestEmptyCreateChatRequestBody(t *testing.T) {
	body := CreateChatRequestBody{}

	actual := body.Validate()
	expected := errors.New("name is empty")

	if !utils.AreErrorsEqual(actual, expected) {
		t.Errorf("Expected: %v, got: %v", actual, expected)
	}
}

func TestConvertValidUpdateChatRequestBodyToMap(t *testing.T) {
	body := UpdateChatRequestBody{
		Name: "New chat's name",
	}

	actual := body.ToMap()
	expected := map[string]any{
		"name": "New chat's name",
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected: %v, got: %v", actual, expected)
	}
}

func TestConvertEmptyUpdateChatRequestBodyToMap(t *testing.T) {
	body := UpdateChatRequestBody{}

	actual := body.ToMap()
	expected := map[string]any{}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected: %v, got: %v", actual, expected)
	}
}

func TestConvertWhiteSpaceUpdateChatRequestBodyToMap(t *testing.T) {
	body := UpdateChatRequestBody{
		Name: "   ",
	}

	actual := body.ToMap()
	expected := map[string]any{}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected: %v, got: %v", actual, expected)
	}
}
