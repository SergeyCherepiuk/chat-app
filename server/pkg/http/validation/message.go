package validation

import (
	"errors"
	"strings"

	"github.com/SergeyCherepiuk/chat-app/domain"
)

type GetHistoryWithNextResponseBody[Message domain.DirectMessage | domain.GroupMessage] struct {
	History []Message `json:"history"`
	Next    string                 `json:"next"`
}

type GetHistoryResponseBody[Message domain.DirectMessage | domain.GroupMessage] struct {
	History []Message `json:"history"`
}

type CreateMessageBody struct {
	Message string `json:"message"`
}

func (body *CreateMessageBody) Validate() error {
	body.Message = strings.TrimSpace(body.Message)

	var err error
	if body.Message == "" {
		err = errors.Join(err, errors.New("message is empty"))
	}
	return err
}

type UpdateMessageRequestBody struct {
	Message string `json:"message"`
}

func (body *UpdateMessageRequestBody) Validate() error {
	body.Message = strings.TrimSpace(body.Message)

	var err error
	if body.Message == "" {
		err = errors.Join(err, errors.New("message is empty"))
	}
	return err
}
