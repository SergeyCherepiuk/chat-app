package storage

import (
	"github.com/SergeyCherepiuk/chat-app/models"
	"gorm.io/gorm"
)

type MessageStorage struct {
	pdb *gorm.DB
}

func NewMessageStorage(pdb *gorm.DB) *MessageStorage {
	return &MessageStorage{pdb: pdb}
}

func (storage MessageStorage) GetAllForChat(chatId uint) ([]models.Message, error) {
	messages := []models.Message{}
	r := storage.pdb.Where("chat_id = ?", chatId).Find(&messages)
	if r.Error != nil {
		return []models.Message{}, r.Error
	}
	return messages, nil
}

func (storage MessageStorage) GetAllForUserInChat(userId uint, chatId uint) ([]models.Message, error) {
	messages := []models.Message{}
	r := storage.pdb.Where("user_id = ? AND chat_id = ?", userId, chatId).Find(&messages)
	if r.Error != nil {
		return []models.Message{}, r.Error
	}
	return messages, nil
}

func (storage MessageStorage) Create(message models.Message) error {
	return storage.pdb.Create(&message).Error
}

func (storage MessageStorage) Update(messageId uint, updatedText string) error {
	message := models.Message{ID: messageId}
	return storage.pdb.Model(&message).Update("message", updatedText).Error
}

func (storage MessageStorage) Delete(messageId uint) error {
	return storage.pdb.Delete(&models.Message{}, messageId).Error
}
