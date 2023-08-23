package postgres

import (
	"github.com/SergeyCherepiuk/chat-app/domain"
	"github.com/SergeyCherepiuk/chat-app/pkg/settings"
	"github.com/SergeyCherepiuk/chat-app/utils"
	"github.com/jmoiron/sqlx"
)

type DirectMessageService struct {
	getHistoryStmt      *sqlx.NamedStmt
	createStmt          *sqlx.NamedStmt
	updateStmt          *sqlx.NamedStmt
	deleteStmt          *sqlx.NamedStmt
	deleteAllStmt       *sqlx.NamedStmt
	isBelongsToChatStmt *sqlx.NamedStmt
	isAuthorStmt        *sqlx.NamedStmt
}

func NewDirectMessageService() *DirectMessageService {
	service := DirectMessageService{}
	utils.MustPrepareNamed(db, &service.getHistoryStmt, `SELECT * FROM direct_messages WHERE ((message_from = :user_id AND message_to = :companion_id) OR (message_from = :companion_id AND message_to = :user_id)) AND id <= :from_id ORDER BY created_at DESC LIMIT :limit`)
	utils.MustPrepareNamed(db, &service.createStmt, `INSERT INTO direct_messages (message_from, message_to, message, is_edited) VALUES (:message_from, :message_to, :message, :is_edited) RETURNING *`)
	utils.MustPrepareNamed(db, &service.updateStmt, `UPDATE direct_messages SET message = :message, is_edited = true WHERE id = :message_id`)
	utils.MustPrepareNamed(db, &service.deleteStmt, `DELETE FROM direct_messages WHERE id = :message_id`)
	utils.MustPrepareNamed(db, &service.deleteAllStmt, `DELETE FROM direct_messages WHERE (message_from = :user_id AND message_to = :companion_id) OR (message_from = :companion_id AND message_to = :user_id)`)
	utils.MustPrepareNamed(db, &service.isBelongsToChatStmt, `SELECT FROM direct_messages WHERE id = :message_id AND ((message_from = :user_id AND message_to = :companion_id) OR (message_from = :companion_id AND message_to = :user_id))`)
	utils.MustPrepareNamed(db, &service.isAuthorStmt, `SELECT FROM direct_messages WHERE id = :message_id AND message_from = :user_id`)
	return &service
}

func (service DirectMessageService) GetHistory(userId, companionId, fromId uint) ([]domain.DirectMessage, error) {
	namedParams := map[string]any{
		"user_id":      userId,
		"companion_id": companionId,
		"from_id":      fromId,
		"limit":        settings.CHAT_HISTORY_BLOCK_SIZE,
	}

	history := []domain.DirectMessage{}
	if err := service.getHistoryStmt.Select(&history, namedParams); err != nil {
		return nil, err
	}

	return history, nil
}

func (service DirectMessageService) Create(message *domain.DirectMessage) error {
	namedParams := map[string]any{
		"message_from": message.From,
		"message_to":   message.To,
		"message":      message.Message,
		"is_edited":    message.IsEdited,
	}

	messageAfterInsert := domain.DirectMessage{}
	if err := service.createStmt.Get(&messageAfterInsert, namedParams); err != nil {
		return err
	}

	message.ID = messageAfterInsert.ID
	message.CreatedAt = messageAfterInsert.CreatedAt
	return nil
}

func (service DirectMessageService) Update(messageId uint, updatedMessage string) error {
	namedParams := map[string]any{
		"message_id": messageId,
		"message":    updatedMessage,
	}

	_, err := service.updateStmt.Exec(namedParams)
	return err
}

func (service DirectMessageService) Delete(messageId uint) error {
	namedParams := map[string]any{
		"message_id": messageId,
	}

	_, err := service.deleteStmt.Exec(namedParams)
	return err
}

func (service DirectMessageService) DeleteAll(userId, companionId uint) error {
	namedParams := map[string]any{
		"user_id":      userId,
		"companion_id": companionId,
	}

	_, err := service.deleteAllStmt.Exec(namedParams)
	return err
}

func (service DirectMessageService) IsBelongsToChat(messageId, userId, companionId uint) (bool, error) {
	namedParams := map[string]any{
		"message_id":   messageId,
		"user_id":      userId,
		"companion_id": companionId,
	}

	if result, err := service.isBelongsToChatStmt.Exec(namedParams); err != nil {
		return false, err
	} else if rowsAffected, err := result.RowsAffected(); err != nil {
		return false, err
	} else {
		return rowsAffected > 0, nil
	}
}

func (service DirectMessageService) IsAuthor(messageId, userId uint) (bool, error) {
	namedParams := map[string]any{
		"message_id": messageId,
		"user_id":    userId,
	}

	if result, err := service.isAuthorStmt.Exec(namedParams); err != nil {
		return false, err
	} else if rowsAffected, err := result.RowsAffected(); err != nil {
		return false, err
	} else {
		return rowsAffected > 0, nil
	}
}
