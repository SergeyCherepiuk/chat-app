package groupchatmessagestorage

import (
	"github.com/SergeyCherepiuk/chat-app/models"
	"github.com/SergeyCherepiuk/chat-app/utils"
	"github.com/jmoiron/sqlx"
)

type GroupChatMessageStorage interface {
	Create(message *models.GroupChatMessage) error
	Update(messageId uint, updates map[string]any) error
	Delete(messageId uint) error
	IsBelongsToChat(messageId, chatId uint) (bool, error)
	IsAuthor(messageId, userId uint) (bool, error)
}

type GroupChatMessageStorageImpl struct {
	pdb                 *sqlx.DB
	createStmt          *sqlx.NamedStmt
	updateStmt          *sqlx.NamedStmt
	deleteStmt          *sqlx.NamedStmt
	isBelongsToChatStmt *sqlx.NamedStmt
	isAuthor            *sqlx.NamedStmt
}

func New(pdb *sqlx.DB) *GroupChatMessageStorageImpl {
	storage := GroupChatMessageStorageImpl{pdb: pdb}
	utils.MustPrepareNamed(pdb, &storage.createStmt, `INSERT INTO group_chat_messages (message, user_id, chat_id, is_edited) VALUES (:message, :user_id, :chat_id, :is_edited) RETURNING *`)
	utils.MustPrepareNamed(pdb, &storage.updateStmt, `UPDATE group_chat_messages SET message = :message WHERE id = :message_id`)
	utils.MustPrepareNamed(pdb, &storage.deleteStmt, `DELETE FROM group_chat_messages WHERE id = :message_id`)
	utils.MustPrepareNamed(pdb, &storage.isBelongsToChatStmt, `SELECT FROM group_chat_messages WHERE id = :message_id AND chat_id = :chat_id`)
	utils.MustPrepareNamed(pdb, &storage.isAuthor, `SELECT FROM group_chat_messages WHERE id = :message_id AND user_id = :user_id`)
	return &GroupChatMessageStorageImpl{pdb: pdb}
}

func (storage GroupChatMessageStorageImpl) Create(message *models.GroupChatMessage) error {
	namedParams := map[string]any{
		"message":   message.Message,
		"user_id":   message.UserID,
		"chat_id":   message.ChatID,
		"is_edited": message.IsEdited,
	}

	insertedMessage := models.GroupChatMessage{}
	if err := storage.createStmt.Get(&insertedMessage, namedParams); err != nil {
		return err
	}

	message.ID = insertedMessage.ID
	message.CreatedAt = insertedMessage.CreatedAt
	return nil
}

func (storage GroupChatMessageStorageImpl) Update(messageId uint, updatedMessage string) error {
	namedParams := map[string]any{
		"message_id": messageId,
		"message":    updatedMessage,
	}

	_, err := storage.updateStmt.Exec(namedParams)
	return err
}

func (storage GroupChatMessageStorageImpl) Delete(messageId uint) error {
	namedParams := map[string]any{
		"message_id": messageId,
	}

	_, err := storage.deleteStmt.Exec(namedParams)
	return err
}

func (storage GroupChatMessageStorageImpl) IsBelongsToChat(messageId, chatId uint) (bool, error) {
	namedParams := map[string]any{
		"message_id": messageId,
		"chat_id":    chatId,
	}

	if result, err := storage.isBelongsToChatStmt.Exec(namedParams); err != nil {
		return false, err
	} else if rowsAffected, err := result.RowsAffected(); err != nil {
		return false, err
	} else {
		return rowsAffected > 0, nil
	}
}

func (storage GroupChatMessageStorageImpl) IsAuthor(messageId, userId uint) (bool, error) {
	namedParams := map[string]any{
		"message_id": messageId,
		"user_id":    userId,
	}

	if result, err := storage.isAuthor.Exec(namedParams); err != nil {
		return false, err
	} else if rowsAffected, err := result.RowsAffected(); err != nil {
		return false, err
	} else {
		return rowsAffected > 0, nil
	}
}
