package chathandler

import "github.com/SergeyCherepiuk/chat-app/storage"

type ChatHandler struct {
	storage *storage.ChatStorage
}

func NewChatHandler(storage *storage.ChatStorage) *ChatHandler {
	return &ChatHandler{storage: storage}
}
