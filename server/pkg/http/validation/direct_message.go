package validation

import (
	"errors"
	"strings"
)

type CreateDirectMessageBody struct {
	Message string `json:"message"`
}

func (body *CreateDirectMessageBody) Validate() error {
	body.Message = strings.TrimSpace(body.Message)

	var err error
	if body.Message == "" {
		err = errors.Join(err, errors.New("message is empty"))
	}
	return err
}

type UpdateDirectMessageRequestBody struct {
	Message string `json:"message"`
}

func (body *UpdateDirectMessageRequestBody) Validate() error {
	body.Message = strings.TrimSpace(body.Message)

	var err error
	if body.Message == "" {
		err = errors.Join(err, errors.New("message is empty"))
	}
	return err
}
