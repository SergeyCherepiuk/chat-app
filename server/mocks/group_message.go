package mocks

import (
	"errors"
	"time"

	"github.com/SergeyCherepiuk/chat-app/domain"
)

type GroupMessageService struct{}

func NewGroupMessageService() *GroupMessageService {
	return &GroupMessageService{}
}

func (service GroupMessageService) reset() {
	groupMessages = []domain.GroupMessage{
		{ID: 1, Message: "First message", UserID: 1, ChatID: 1, IsEdited: false, CreatedAt: time.Now()},
		{ID: 1, Message: "Second message", UserID: 1, ChatID: 2, IsEdited: false, CreatedAt: time.Now()},
		{ID: 1, Message: "Third message", UserID: 1, ChatID: 3, IsEdited: true, CreatedAt: time.Now()},
	}
}

func (service GroupMessageService) Create(message *domain.GroupMessage) error {
	service.reset()
	groupMessages = append(groupMessages, *message)
	return nil
}

func (service GroupMessageService) Update(messageId uint, updates map[string]any) error {
	service.reset()
	for _, message := range groupMessages {
		if message.ID == messageId {
			return nil
		}
	}
	return errors.New("message not found")
}

func (service GroupMessageService) Delete(messageId uint) error {
	service.reset()
	for i, message := range groupMessages {
		if message.ID == messageId {
			groupMessages = append(groupMessages[:i], groupMessages[i+1:]...)
			return nil
		}
	}
	return errors.New("message not found")
}

func (service GroupMessageService) IsBelongsToChat(messageId, chatId uint) (bool, error) {
	service.reset()
	for _, message := range groupMessages {
		if message.ID == messageId {
			return message.ChatID == chatId, nil
		}
	}
	return false, errors.New("message not found")
}

func (service GroupMessageService) IsAuthor(messageId, userId uint) (bool, error) {
	service.reset()
	for _, message := range groupMessages {
		if message.ID == messageId {
			return message.UserID == userId, nil
		}
	}
	return false, errors.New("message not found")
}
