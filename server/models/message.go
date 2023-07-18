package models

type Message struct {
	ID      uint   `json:"id" db:"id" gorm:"primary"`
	Message string `json:"message" db:"message"`
	UserID  uint   `json:"user_id" db:"user_id"`
	ChatID  uint   `json:"chat_id" db:"chat_id"`
}
