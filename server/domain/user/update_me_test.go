package userdomain_test

import (
	"reflect"
	"testing"

	userdomain "github.com/SergeyCherepiuk/chat-app/domain/user"
)

func TestInvalidUpdateUserRequestBody(t *testing.T) {
	body := userdomain.UpdateUserRequestBody{
		FirstName:   "",
		LastName:    " ",
		Username:    "  ",
		Description: "   ",
	}

	actualMap := body.ToMap()
	expectedMap := map[string]any{}
	expectedBody := userdomain.UpdateUserRequestBody{
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

func TestPartiallyValidUpdateUserRequestBody(t *testing.T) {
	body := userdomain.UpdateUserRequestBody{
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
	expectedBody := userdomain.UpdateUserRequestBody{
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

func TestValidUpdateUserRequestBody(t *testing.T) {
	body := userdomain.UpdateUserRequestBody{
		FirstName:   "John",
		LastName:    "Doe",
		Username:    "johndoe",
		Description: "New John's description",
	}

	actualMap := body.ToMap()
	expectedMap := map[string]any{
		"first_name": "John",
		"last_name": "Doe",
		"username":    "johndoe",
		"description": "New John's description",
	}
	expectedBody := userdomain.UpdateUserRequestBody{
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