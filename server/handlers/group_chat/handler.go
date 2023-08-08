package groupchathandler

import (
	groupchatstorage "github.com/SergeyCherepiuk/chat-app/storage/group_chat"
)

type GroupChatHandler struct {
	storage groupchatstorage.GroupChatStorage
}

func New(storage groupchatstorage.GroupChatStorage) *GroupChatHandler {
	return &GroupChatHandler{storage: storage}
}
