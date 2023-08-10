package validation

import (
	"errors"
	"strings"
)

type GetGroupChatResponseBody struct {
	Name string `json:"name"`
}

type CreateGroupChatRequestBody struct {
	Name string `json:"name"`
}

func (body *CreateGroupChatRequestBody) Validate() error {
	body.Name = strings.TrimSpace(body.Name)

	var err error
	if body.Name == "" {
		err = errors.Join(err, errors.New("chat name is empty"))
	}

	return err
}

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
