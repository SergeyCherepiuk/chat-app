package chathandler

import chatstorage "github.com/SergeyCherepiuk/chat-app/storage/chat"

type ChatHandler struct {
	storage chatstorage.ChatStorage
}

func New(storage chatstorage.ChatStorage) *ChatHandler {
	return &ChatHandler{storage: storage}
}
