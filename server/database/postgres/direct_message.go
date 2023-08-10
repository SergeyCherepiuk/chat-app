package postgres

import (
	"github.com/SergeyCherepiuk/chat-app/domain"
)

type DirectMessageService struct{}

func NewDirectMessageService() *DirectMessageService {
	return &DirectMessageService{}
}

func (service DirectMessageService) GetHistory(userId, companionId uint) ([]domain.DirectMessage, error) {
	query := `SELECT * FROM direct_messages WHERE (message_from = :user_id AND message_to = :companion_id) OR (message_from = :companion_id AND message_to = :user_id) ORDER BY created_at`
	namedParams := map[string]any{
		"user_id":      userId,
		"companion_id": companionId,
	}

	stmt, err := db.PrepareNamed(query)
	if err != nil {
		return []domain.DirectMessage{}, err
	}
	defer stmt.Close()

	history := []domain.DirectMessage{}
	if err := stmt.Select(&history, namedParams); err != nil {
		return []domain.DirectMessage{}, err
	}

	return history, nil
}

func (service DirectMessageService) Create(message *domain.DirectMessage) error {
	query := `INSERT INTO direct_messages (message_from, message_to, message, is_edited) VALUES (:message_from, :message_to, :message, :is_edited) RETURNING *`
	namedParams := map[string]any{
		"message_from": message.From,
		"message_to":   message.To,
		"message":      message.Message,
		"is_edited":    message.IsEdited,
	}

	stmt, err := db.PrepareNamed(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	messageAfterInsert := domain.DirectMessage{}
	if err := stmt.Get(&messageAfterInsert, namedParams); err != nil {
		return err
	}

	message.ID = messageAfterInsert.ID
	message.CreatedAt = messageAfterInsert.CreatedAt
	return nil
}

func (service DirectMessageService) Update(messageId uint, updatedMessage string) error {
	query := `UPDATE direct_messages SET message = :message, is_edited = true WHERE id = :message_id`
	namedParams := map[string]any{
		"message_id": messageId,
		"message":    updatedMessage,
	}

	stmt, err := db.PrepareNamed(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(namedParams)
	return err
}

func (service DirectMessageService) Delete(messageId uint) error {
	query := `DELETE FROM direct_messages WHERE id = :message_id`
	namedParams := map[string]any{
		"message_id": messageId,
	}

	stmt, err := db.PrepareNamed(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(namedParams)
	return err
}

func (service DirectMessageService) DeleteAll(userId, companionId uint) error {
	query := `DELETE FROM direct_messages WHERE (message_from = :user_id AND message_to = :companion_id) OR (message_from = :companion_id AND message_to = :user_id)`
	namedParams := map[string]any{
		"user_id":      userId,
		"companion_id": companionId,
	}

	stmt, err := db.PrepareNamed(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(namedParams)
	return err
}

func (service DirectMessageService) IsBelongsToChat(messageId, userId, companionId uint) (bool, error) {
	query := `SELECT FROM direct_messages WHERE id = :message_id AND ((message_from = :user_id AND message_to = :companion_id) OR (message_from = :companion_id AND message_to = :user_id))`
	namedParams := map[string]any{
		"message_id":   messageId,
		"user_id":      userId,
		"companion_id": companionId,
	}

	stmt, err := db.PrepareNamed(query)
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

func (service DirectMessageService) IsAuthor(messageId, userId uint) (bool, error) {
	query := `SELECT FROM direct_messages WHERE id = :message_id AND message_from = :user_id`
	namedParams := map[string]any{
		"message_id": messageId,
		"user_id":    userId,
	}

	stmt, err := db.PrepareNamed(query)
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
