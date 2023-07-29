package userhandler_test

import (
	"reflect"
	"testing"

	userhandler "github.com/SergeyCherepiuk/chat-app/handlers/user"
)

func TestConvertValidUpdateUserRequestBodyToMap(t *testing.T) {
	body := userhandler.UpdateUserRequestBody{
		FirstName: "Andrew",
		LastName:  "Brown",
		Username:  "andrewbrown",
	}

	actual := body.ToMap()
	expected := map[string]any{
		"first_name": "Andrew",
		"last_name":  "Brown",
		"username":   "andrewbrown",
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected: %v, got: %v", actual, expected)
	}
}

func TestConvertEmptyUpdateUserRequestBodyToMap(t *testing.T) {
	body := userhandler.UpdateUserRequestBody{}

	actual := body.ToMap()
	expected := map[string]any{}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected: %v, got: %v", actual, expected)
	}
}

func TestConvertWhiteSpaceUpdateUserRequestBodyToMap(t *testing.T) {
	body := userhandler.UpdateUserRequestBody{
		FirstName: "",
		LastName:  " ",
		Username:  "  ",
	}

	actual := body.ToMap()
	expected := map[string]any{}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected: %v, got: %v", actual, expected)
	}
}
