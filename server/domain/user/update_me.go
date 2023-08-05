package userdomain

import "strings"

type UpdateUserRequestBody struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
}

func (body *UpdateUserRequestBody) ToMap() map[string]any {
	body.FirstName = strings.TrimSpace(body.FirstName)
	body.LastName = strings.TrimSpace(body.LastName)
	body.Username = strings.TrimSpace(body.Username)

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

	return updates
}
