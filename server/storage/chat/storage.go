package chatstorage

import (
	"github.com/SergeyCherepiuk/chat-app/models"
	"github.com/jmoiron/sqlx"
)

type ChatStorage interface {
	GetHistory(userId, companionId uint) ([]models.ChatMessage, error)
	Create(message *models.ChatMessage) error
	Update(messageId uint, updatedMessage string) error
	Delete(messageId uint) error
	DeleteAll(userId, companionId uint) error
	IsBelongsToChat(messageId, userId, companionId uint) (bool, error)
	IsAuthor(messageId, userId uint) (bool, error)
}

type ChatStorageImpl struct {
	pdb *sqlx.DB
}

func New(pdb *sqlx.DB) *ChatStorageImpl {
	return &ChatStorageImpl{pdb: pdb}
}

func (storage ChatStorageImpl) GetHistory(userId, companionId uint) ([]models.ChatMessage, error) {
	query := `SELECT * FROM chat_messages WHERE (message_from = :user_id AND message_to = :companion_id) OR (message_from = :companion_id AND message_to = :user_id) ORDER BY created_at`
	namedParams := map[string]any{
		"user_id":      userId,
		"companion_id": companionId,
	}

	stmt, err := storage.pdb.PrepareNamed(query)
	if err != nil {
		return []models.ChatMessage{}, err
	}
	defer stmt.Close()

	history := []models.ChatMessage{}
	if err := stmt.Select(&history, namedParams); err != nil {
		return []models.ChatMessage{}, err
	}

	return messages, nil
}

func (storage ChatStorageImpl) Create(message *models.ChatMessage) error {
	query := `INSERT INTO chat_messages (message_from, message_to, message, is_edited) VALUES (:message_from, :message_to, :message, :is_edited) RETURNING *`
	namedParams := map[string]any{
		"message_from": message.From,
		"message_to":   message.To,
		"message":      message.Message,
		"is_edited":    message.IsEdited,
	}

	stmt, err := storage.pdb.PrepareNamed(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	messageAfterInsert := models.ChatMessage{}
	if err := stmt.Get(&messageAfterInsert, namedParams); err != nil {
		return err
	}

	message.ID = messageAfterInsert.ID
	message.CreatedAt = messageAfterInsert.CreatedAt
	return nil
}

func (storage ChatStorageImpl) Update(messageId uint, updatedMessage string) error {
	query := `UPDATE chat_messages SET message = :message, is_edited = true WHERE id = :message_id`
	namedParams := map[string]any{
		"message_id": messageId,
		"message":    updatedMessage,
	}

	stmt, err := storage.pdb.PrepareNamed(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(namedParams)
	return err
}

func (storage ChatStorageImpl) Delete(messageId uint) error {
	query := `DELETE FROM chat_messages WHERE id = :message_id`
	namedParams := map[string]any{
		"message_id": messageId,
	}

	stmt, err := storage.pdb.PrepareNamed(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(namedParams)
	return err
}

func (storage ChatStorageImpl) DeleteAll(userId, companionId uint) error {
	query := `DELETE FROM chat_messages WHERE (message_from = :user_id AND message_to = :companion_id) OR (message_from = :companion_id AND message_to = :user_id)`
	namedParams := map[string]any{
		"user_id":      userId,
		"companion_id": companionId,
	}

	stmt, err := storage.pdb.PrepareNamed(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(namedParams)
	return err
}

func (storage ChatStorageImpl) IsBelongsToChat(messageId, userId, companionId uint) (bool, error) {
	query := `SELECT FROM chat_messages WHERE id = :message_id AND ((message_from = :user_id AND message_to = :companion_id) OR (message_from = :companion_id AND message_to = :user_id))`
	namedParams := map[string]any{
		"message_id":   messageId,
		"user_id":      userId,
		"companion_id": companionId,
	}

	stmt, err := storage.pdb.PrepareNamed(query)
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	if result, err := stmt.Exec(namedParams); err != nil {
		return false, err
	} else if rowsAffected, err := result.RowsAffected(); err != nil {
		return false, err
	} else {
		return rowsAffected > 0, nil
	}
}

func (storage ChatStorageImpl) IsAuthor(messageId, userId uint) (bool, error) {
	query := `SELECT FROM chat_messages WHERE id = :message_id AND message_from = :user_id`
	namedParams := map[string]any{
		"message_id": messageId,
		"user_id":    userId,
	}

	stmt, err := storage.pdb.PrepareNamed(query)
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	if result, err := stmt.Exec(namedParams); err != nil {
		return false, err
	} else if rowsAffected, err := result.RowsAffected(); err != nil {
		return false, err
	} else {
		return rowsAffected > 0, nil
	}
}
