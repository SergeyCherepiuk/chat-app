package domain

import (
	"errors"
	"strings"
)

type CreateChatRequestBody struct {
	Name string `json:"name"`
}

func (body CreateChatRequestBody) Validate() error {
	var err error

	if strings.TrimSpace(body.Name) == "" {
		err = errors.Join(err, errors.New("name is empty"))
	}

	return err
}

type UpdateChatRequestBody struct {
	Name string `json:"name"`
}

func (body UpdateChatRequestBody) ToMap() map[string]any {
	updates := make(map[string]any)

	if strings.TrimSpace(body.Name) != "" {
		updates["name"] = body.Name
	}

	return updates
}
