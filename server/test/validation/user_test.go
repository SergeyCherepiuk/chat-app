package validation_test

import (
	"reflect"
	"testing"

	"github.com/SergeyCherepiuk/chat-app/pkg/http/validation"
)

func Test_Invalid_UpdateUserRequestBody(t *testing.T) {
	body := validation.UpdateUserRequestBody{
		FirstName:   "",
		LastName:    " ",
		Username:    "  ",
		Description: "   ",
	}

	actualMap := body.ToMap()
	expectedMap := map[string]any{}
	expectedBody := validation.UpdateUserRequestBody{
		FirstName:   "",
		LastName:    "",
		Username:    "",
		Description: "",
	}

	if !reflect.DeepEqual(actualMap, expectedMap) {
		t.Errorf("expected: %v, actual: %v\n", expectedMap, actualMap)
	}
	if body != expectedBody {
		t.Errorf("expected: %v, actual: %v\n", expectedBody, body)
	}
}

func Test_PartiallyValid_UpdateUserRequestBody(t *testing.T) {
	body := validation.UpdateUserRequestBody{
		FirstName:   "",
		LastName:    "   ",
		Username:    "newusername",
		Description: "New description   ",
	}

	actualMap := body.ToMap()
	expectedMap := map[string]any{
		"username":    "newusername",
		"description": "New description",
	}
	expectedBody := validation.UpdateUserRequestBody{
		FirstName:   "",
		LastName:    "",
		Username:    "newusername",
		Description: "New description",
	}

	if !reflect.DeepEqual(actualMap, expectedMap) {
		t.Errorf("expected: %v, actual: %v\n", expectedMap, actualMap)
	}
	if body != expectedBody {
		t.Errorf("expected: %v, actual: %v\n", expectedBody, body)
	}
}

func Test_Valid_UpdateUserRequestBody(t *testing.T) {
	body := validation.UpdateUserRequestBody{
		FirstName:   "John",
		LastName:    "Doe",
		Username:    "johndoe",
		Description: "New John's description",
	}

	actualMap := body.ToMap()
	expectedMap := map[string]any{
		"first_name":  "John",
		"last_name":   "Doe",
		"username":    "johndoe",
		"description": "New John's description",
	}
	expectedBody := validation.UpdateUserRequestBody{
		FirstName:   "John",
		LastName:    "Doe",
		Username:    "johndoe",
		Description: "New John's description",
	}

	if !reflect.DeepEqual(actualMap, expectedMap) {
		t.Errorf("expected: %v, actual: %v\n", expectedMap, actualMap)
	}
	if body != expectedBody {
		t.Errorf("expected: %v, actual: %v\n", expectedBody, body)
	}
}
