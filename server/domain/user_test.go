package domain

import (
	"reflect"
	"testing"
)

func TestConvertValidUpdateUserRequestBodyToMap(t *testing.T) {
	body := UpdateUserRequestBody{
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
		t.Errorf("Expected: %v, got: %v", actual, expected)
	}
}

func TestConvertEmptyUpdateUserRequestBodyToMap(t *testing.T) {
	body := UpdateUserRequestBody{}

	actual := body.ToMap()
	expected := map[string]any{}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected: %v, got: %v", actual, expected)
	}
}

func TestConvertWhiteSpaceUpdateUserRequestBodyToMap(t *testing.T) {
	body := UpdateUserRequestBody{
		FirstName: "",
		LastName:  " ",
		Username:  "  ",
	}

	actual := body.ToMap()
	expected := map[string]any{}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected: %v, got: %v", actual, expected)
	}
}
