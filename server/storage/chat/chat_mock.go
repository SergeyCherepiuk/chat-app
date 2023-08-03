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

var chats = []models.Chat{
	{ID: 1, Name: "First chat", Messages: []models.Message{messages[0], messages[1]}},
}

var messages = []models.Message{
	{ID: 1, Message: "First message", SentAt: time.Now(), UserID: 1, ChatID: 1},
	{ID: 2, Message: "Second message", SentAt: time.Now(), UserID: 2, ChatID: 1},
}

func (storage ChatStorageMock) GetAllChats() ([]models.Chat, error) {
	return chats, nil
}

func (storage ChatStorageMock) GetChatById(chatId uint) (models.Chat, error) {
	for _, c := range chats {
		if c.ID == chatId {
			return c, nil
		}
	}
	return models.Chat{}, errors.New("chat not found")
}

func (storage ChatStorageMock) CreateChat(chat models.Chat) error {
	return nil
}

func (storage ChatStorageMock) UpdateChat(chatId uint, updates map[string]any) error {
	_, err := storage.GetChatById(chatId)
	return err
}

func (storage ChatStorageMock) DeleteChat(chatId uint) error {
	_, err := storage.GetChatById(chatId)
	return err
}

func (storage ChatStorageMock) IsChatExists(chatId uint) bool {
	for _, c := range chats {
		if c.ID == chatId {
			return true
		}
	}
	return false
}

func (storage ChatStorageMock) GetAllMessages(chatId uint) ([]models.Message, error) {
	return messages, nil
}

func (storage ChatStorageMock) CreateMessage(message *models.Message) error {
	return nil
}

func (storage ChatStorageMock) UpdateMessage(messageId uint, updatedMessage string) error {
	for _, m := range messages {
		if m.ID == messageId {
			return nil
		}
	}
	return errors.New("message not found")
}

func (storage ChatStorageMock) DeleteMessage(messageId uint) error {
	for _, m := range messages {
		if m.ID == messageId {
			return nil
		}
	}
	return errors.New("message not found")
}
