package groupchatstorage

import (
	"errors"
	"time"

	"github.com/SergeyCherepiuk/chat-app/models"
)

type GroupChatStorageMock struct{}

func NewMock() *GroupChatStorageMock {
	return &GroupChatStorageMock{}
}

var chats []models.GroupChat
var messages []models.GroupChatMessage

func (storage GroupChatStorageMock) reset() {
	chats = []models.GroupChat{
		{ID: 1, Name: "First group chat", CreatedAt: time.Now()},
	}
	messages = []models.GroupChatMessage{
		{ID: 1, Message: "First message", UserID: 1, ChatID: 1, IsEdited: false, CreatedAt: time.Now()},
		{ID: 1, Message: "Second message", UserID: 1, ChatID: 2, IsEdited: false, CreatedAt: time.Now()},
		{ID: 1, Message: "Third message", UserID: 1, ChatID: 3, IsEdited: true, CreatedAt: time.Now()},
	}
}

func (storage GroupChatStorageMock) GetChat(chatId uint) (models.GroupChat, error) {
	storage.reset()
	for _, chat := range chats {
		if chat.ID == chatId {
			return chat, nil
		}
	}
	return models.GroupChat{}, errors.New("chat not found")
}

func (storage GroupChatStorageMock) GetHistory(chatId uint) ([]models.GroupChatMessage, error) {
	storage.reset()
	history := []models.GroupChatMessage{}
	for _, message := range messages {
		if message.ChatID == chatId {
			history = append(history, message)
		}
	}
	return history, nil
}

func (storage GroupChatStorageMock) Create(chat *models.GroupChat) error {
	storage.reset()
	chats = append(chats, *chat)
	return nil
}

func (storage GroupChatStorageMock) Update(chatId uint, updates map[string]any) error {
	storage.reset()
	_, err := storage.GetChat(chatId)
	return err
}

func (storage GroupChatStorageMock) Delete(chatId uint) error {
	storage.reset()
	for i, chat := range chats {
		if chat.ID == chatId {
			chats = append(chats[:i], chats[i+1:]...)
			return nil
		}
	}
	return errors.New("chat not found")
}

func (storage GroupChatStorageMock) IsAdmin(chatId, userId uint) (bool, error) {
	storage.reset()
	for _, chat := range chats {
		if chat.ID == chatId {
			return chat.CreatorID == userId, nil
		}
	}
	return false, errors.New("chat not found")
}
