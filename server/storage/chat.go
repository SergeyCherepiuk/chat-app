package storage

import (
	"github.com/SergeyCherepiuk/chat-app/models"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type ChatStorage struct {
	pdb *gorm.DB
	rdb *redis.Client
}

func NewChatStorage(pdb *gorm.DB) *ChatStorage {
	return &ChatStorage{pdb: pdb}
}

func (storage ChatStorage) IsExists(chatId uint) bool {
	r := storage.pdb.First(&models.Chat{}, chatId)
	return r.Error == nil && r.RowsAffected > 0
}

func (storage ChatStorage) GetAllForChat(chatId uint) ([]models.Message, error) {
	messages := []models.Message{}
	r := storage.pdb.Where("chat_id = ?", chatId).Find(&messages)
	if r.Error != nil {
		return []models.Message{}, r.Error
	}
	return messages, nil
}

func (storage ChatStorage) Create(message models.Message) error {
	return storage.pdb.Create(&message).Error
}

func (storage ChatStorage) Update(messageId uint, updatedText string) error {
	message := models.Message{ID: messageId}
	return storage.pdb.Model(&message).Update("message", updatedText).Error
}

func (storage ChatStorage) Delete(messageId uint) error {
	return storage.pdb.Delete(&models.Message{}, messageId).Error
}
