package validation

import "strings"

type GetUserResponseBody struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Username    string `json:"username"`
	Description string `json:"description"`
}

type UpdateUserRequestBody struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Username    string `json:"username"`
	Description string `json:"description"`
}

func (body *UpdateUserRequestBody) ToMap() map[string]any {
	body.FirstName = strings.TrimSpace(body.FirstName)
	body.LastName = strings.TrimSpace(body.LastName)
	body.Username = strings.TrimSpace(body.Username)
	body.Description = strings.TrimSpace(body.Description)

	updates := make(map[string]any)
	if body.FirstName != "" {
		updates["first_name"] = body.FirstName
	}
	if body.LastName != "" {
		updates["last_name"] = body.LastName
	}
	if body.Username != "" {
		updates["username"] = body.Username
	}
	if body.Description != "" {
		updates["description"] = body.Description
	}

	return updates
}
