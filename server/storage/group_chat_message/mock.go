package groupchatmessagestorage

import (
	"errors"
	"time"

	"github.com/SergeyCherepiuk/chat-app/models"
)

type GroupChatMessageStorageMock struct{}

func NewMock() *GroupChatMessageStorageMock {
	return &GroupChatMessageStorageMock{}
}

var messages []models.GroupChatMessage

func (storage GroupChatMessageStorageMock) reset() {
	messages = []models.GroupChatMessage{
		{ID: 1, Message: "First message", UserID: 1, ChatID: 1, IsEdited: false, CreatedAt: time.Now()},
		{ID: 1, Message: "Second message", UserID: 1, ChatID: 2, IsEdited: false, CreatedAt: time.Now()},
		{ID: 1, Message: "Third message", UserID: 1, ChatID: 3, IsEdited: true, CreatedAt: time.Now()},
	}
}

func (storage GroupChatMessageStorageMock) Create(message *models.GroupChatMessage) error {
	storage.reset()
	messages = append(messages, *message)
	return nil
}

func (storage GroupChatMessageStorageMock) Update(messageId uint, updates map[string]any) error {
	storage.reset()
	for _, message := range messages {
		if message.ID == messageId {
			return nil
		}
	}
	return errors.New("message not found")
}

func (storage GroupChatMessageStorageMock) Delete(messageId uint) error {
	storage.reset()
	for i, message := range messages {
		if message.ID == messageId {
			messages = append(messages[:i], messages[i+1:]...)
			return nil
		}
	}
	return errors.New("message not found")
}

func (storage GroupChatMessageStorageMock) IsBelongsToChat(messageId, chatId uint) (bool, error) {
	storage.reset()
	for _, message := range messages {
		if message.ID == messageId {
			return message.ChatID == chatId, nil
		}
	}
	return false, errors.New("message not found")
}

func (storage GroupChatMessageStorageMock) IsAuthor(messageId, userId uint) (bool, error) {
	storage.reset()
	for _, message := range messages {
		if message.ID == messageId {
			return message.UserID == userId, nil
		}
	}
	return false, errors.New("message not found")
}
