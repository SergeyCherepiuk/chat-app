package groupchatmessagehandler

import groupchatmessagestorage "github.com/SergeyCherepiuk/chat-app/storage/group_chat_message"

type GroupChatMessageHandler struct {
	storage groupchatmessagestorage.GroupChatMessageStorage
}

func New(storage groupchatmessagestorage.GroupChatMessageStorage) *GroupChatMessageHandler {
	return &GroupChatMessageHandler{storage: storage}
}