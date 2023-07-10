package domain

import (
	"strings"
)

type GetMeResponseBody struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
}

type UpdateMeRequestBody struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
}

func (body UpdateMeRequestBody) ToMap() map[string]any {
	updates := make(map[string]any)
	if strings.TrimSpace(body.FirstName) != "" {
		updates["first_name"] = body.FirstName
	}
	if strings.TrimSpace(body.LastName) != "" {
		updates["last_name"] = body.LastName
	}
	if strings.TrimSpace(body.Username) != "" {
		updates["username"] = body.Username
	}
	return updates
}
