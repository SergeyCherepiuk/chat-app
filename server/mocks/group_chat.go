package mocks

import (
	"errors"
	"time"

	"github.com/SergeyCherepiuk/chat-app/domain"
)

type GroupChatService struct{}

func NewGroupChatService() *GroupChatService {
	return &GroupChatService{}
}

func (service GroupChatService) reset() {
	groupChats = []domain.GroupChat{
		{ID: 1, Name: "First group chat", CreatedAt: time.Now()},
	}
	groupMessages = []domain.GroupMessage{
		{ID: 1, Message: "First message", UserID: 1, ChatID: 1, IsEdited: false, CreatedAt: time.Now()},
		{ID: 1, Message: "Second message", UserID: 1, ChatID: 2, IsEdited: false, CreatedAt: time.Now()},
		{ID: 1, Message: "Third message", UserID: 1, ChatID: 3, IsEdited: true, CreatedAt: time.Now()},
	}
}

func (service GroupChatService) GetChat(chatId uint) (domain.GroupChat, error) {
	service.reset()
	for _, chat := range groupChats {
		if chat.ID == chatId {
			return chat, nil
		}
	}
	return domain.GroupChat{}, errors.New("chat not found")
}

func (service GroupChatService) GetHistory(chatId uint) ([]domain.GroupMessage, error) {
	service.reset()
	history := []domain.GroupMessage{}
	for _, message := range groupMessages {
		if message.ChatID == chatId {
			history = append(history, message)
		}
	}
	return history, nil
}

func (service GroupChatService) Create(chat *domain.GroupChat) error {
	service.reset()
	groupChats = append(groupChats, *chat)
	return nil
}

func (service GroupChatService) Update(chatId uint, updates map[string]any) error {
	service.reset()
	_, err := service.GetChat(chatId)
	return err
}

func (service GroupChatService) Delete(chatId uint) error {
	service.reset()
	for i, chat := range groupChats {
		if chat.ID == chatId {
			groupChats = append(groupChats[:i], groupChats[i+1:]...)
			return nil
		}
	}
	return errors.New("chat not found")
}

func (service GroupChatService) IsAdmin(chatId, userId uint) (bool, error) {
	service.reset()
	for _, chat := range groupChats {
		if chat.ID == chatId {
			return chat.CreatorID == userId, nil
		}
	}
	return false, errors.New("chat not found")
}
