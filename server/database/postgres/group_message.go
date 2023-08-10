package postgres

import (
	"github.com/SergeyCherepiuk/chat-app/domain"
	"github.com/SergeyCherepiuk/chat-app/utils"
	"github.com/jmoiron/sqlx"
)

type GroupMessageService struct {
	createStmt          *sqlx.NamedStmt
	updateStmt          *sqlx.NamedStmt
	deleteStmt          *sqlx.NamedStmt
	isBelongsToChatStmt *sqlx.NamedStmt
	isAuthor            *sqlx.NamedStmt
}

func NewGroupMessageServiceMessageService() *GroupMessageService {
	service := GroupMessageService{}
	utils.MustPrepareNamed(db, &service.createStmt, `INSERT INTO group_messages (message, user_id, chat_id, is_edited) VALUES (:message, :user_id, :chat_id, :is_edited) RETURNING *`)
	utils.MustPrepareNamed(db, &service.updateStmt, `UPDATE group_messages SET message = :message WHERE id = :message_id`)
	utils.MustPrepareNamed(db, &service.deleteStmt, `DELETE FROM group_messages WHERE id = :message_id`)
	utils.MustPrepareNamed(db, &service.isBelongsToChatStmt, `SELECT FROM group_messages WHERE id = :message_id AND chat_id = :chat_id`)
	utils.MustPrepareNamed(db, &service.isAuthor, `SELECT FROM group_messages WHERE id = :message_id AND user_id = :user_id`)
	return &service
}

func (service GroupMessageService) Create(message *domain.GroupMessage) error {
	namedParams := map[string]any{
		"message":   message.Message,
		"user_id":   message.UserID,
		"chat_id":   message.ChatID,
		"is_edited": message.IsEdited,
	}

	insertedMessage := domain.GroupMessage{}
	if err := service.createStmt.Get(&insertedMessage, namedParams); err != nil {
		return err
	}

	message.ID = insertedMessage.ID
	message.CreatedAt = insertedMessage.CreatedAt
	return nil
}

func (service GroupMessageService) Update(messageId uint, updatedMessage string) error {
	namedParams := map[string]any{
		"message_id": messageId,
		"message":    updatedMessage,
	}

	_, err := service.updateStmt.Exec(namedParams)
	return err
}

func (service GroupMessageService) Delete(messageId uint) error {
	namedParams := map[string]any{
		"message_id": messageId,
	}

	_, err := service.deleteStmt.Exec(namedParams)
	return err
}

func (service GroupMessageService) IsBelongsToChat(messageId, chatId uint) (bool, error) {
	namedParams := map[string]any{
		"message_id": messageId,
		"chat_id":    chatId,
	}

	if result, err := service.isBelongsToChatStmt.Exec(namedParams); err != nil {
		return false, err
	} else if rowsAffected, err := result.RowsAffected(); err != nil {
		return false, err
	} else {
		return rowsAffected > 0, nil
	}
}

func (service GroupMessageService) IsAuthor(messageId, userId uint) (bool, error) {
	namedParams := map[string]any{
		"message_id": messageId,
		"user_id":    userId,
	}

	if result, err := service.isAuthor.Exec(namedParams); err != nil {
		return false, err
	} else if rowsAffected, err := result.RowsAffected(); err != nil {
		return false, err
	} else {
		return rowsAffected > 0, nil
	}
}
