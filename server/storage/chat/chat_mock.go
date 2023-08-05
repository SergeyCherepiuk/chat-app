package chatstorage

import (
	"errors"
	"time"

	"github.com/SergeyCherepiuk/chat-app/models"
)

type ChatStorageMock struct{}

func NewMock() *ChatStorageMock {
	return &ChatStorageMock{}
}

var messages = []models.ChatMessage{
	{ID: 1, Message: "First message", From: 1, To: 2, CreatedAt: time.Now()},
	{ID: 2, Message: "Second message", From: 2, To: 1, CreatedAt: time.Now()},
}

func (storage ChatStorageMock) GetChatHistory(userId, companionId uint) ([]models.ChatMessage, error) {
	history := []models.ChatMessage{}
	for _, message := range messages {
		if (message.From == userId && message.To == companionId) || (message.From == companionId && message.To == userId) {
			history = append(history, message)
		}
	}
	return history, nil
}

func (storage ChatStorageMock) DeleteChat(userId, companionId uint) error {
	for i, message := range messages {
		if (message.From == userId && message.To == companionId) || (message.From == companionId && message.To == userId) {
			messages = append(messages[:i], messages[i+1:]...)
		}
	}
	return nil
}

func (storage ChatStorageMock) CreateMessage(message *models.ChatMessage) error {
	messages = append(messages, *message)
	return nil
}

func (storage ChatStorageMock) IsMessageBelongToChat(messageId, userId, companionId uint) (bool, error) {
	for _, message := range messages {
		if message.ID == messageId && ((message.From == userId && message.To == companionId) || (message.From == companionId && message.To == userId)) {
			return true, nil
		}
	}
	return false, errors.New("message not found in chat")
}

func (storage ChatStorageMock) IsAuthor(messageId, userId uint) (bool, error) {
	for _, message := range messages {
		if message.ID == messageId {
			return message.From == userId, nil
		}
	}
	return false, errors.New("message not found")
}

func (storage ChatStorageMock) UpdateMessage(messageId uint, updatedMessage string) error {
	for _, message := range messages {
		if message.ID == messageId {
			message.Message = updatedMessage
			message.IsEdited = true
			return nil
		}
	}
	return errors.New("message not found")
}

func (storage ChatStorageMock) DeleteMessage(messageId uint) error {
	for i, message := range messages {
		if message.ID == messageId {
			messages = append(messages[:i], messages[i+1:]...)
			return nil
		}
	}
	return errors.New("message not found")
}
