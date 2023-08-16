package validation

import "github.com/SergeyCherepiuk/chat-app/domain"

type GetHistoryWithNextResponseBody struct {
	History []domain.DirectMessage `json:"history"`
	Next    string                 `json:"next"`
}

type GetHistoryResponseBody struct {
	History []domain.DirectMessage `json:"history"`
}
