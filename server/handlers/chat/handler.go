package chathandler

import (
	chatstorage "github.com/SergeyCherepiuk/chat-app/storage/chat"
	userstorage "github.com/SergeyCherepiuk/chat-app/storage/user"
)

type ChatHandler struct {
	chatStorage chatstorage.ChatStorage
	userStorage userstorage.UserStorage
}

func New(chatStorage chatstorage.ChatStorage, userStorage userstorage.UserStorage) *ChatHandler {
	return &ChatHandler{chatStorage: chatStorage, userStorage: userStorage}
}
