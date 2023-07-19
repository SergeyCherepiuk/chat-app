package domain

import (
	"errors"
	"strings"
)

type CreateMessageRequestBody struct {
	Text string `json:"text"`
}

func (body CreateMessageRequestBody) Validate() error {
	var err error
	if strings.TrimSpace(body.Text) == "" {
		err = errors.Join(err, errors.New("message body is empty"))
	}
	return err
}
