package groupchatdomain

import (
	"errors"
	"strings"
)

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