package groupchatdomain

import (
	"strings"
)

type UpdateGroupChatRequestBody struct {
	Name string `json:"name"`
}

func (body *UpdateGroupChatRequestBody) ToMap() map[string]any {
	body.Name = strings.TrimSpace(body.Name)

	updates := make(map[string]any)
	if body.Name != "" {
		updates["name"] = body.Name
	}

	return updates
}
