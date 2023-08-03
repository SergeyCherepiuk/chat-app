package models

import "time"

type Message struct {
	ID      uint      `json:"id" db:"id" gorm:"primary"`
	Message string    `json:"message" db:"message" gorm:"not null"`
	SentAt  time.Time `json:"sent_at" db:"sent_at" gorm:"not null"`
	UserID  uint      `json:"user_id" db:"user_id"`
	ChatID  uint      `json:"chat_id" db:"chat_id" gorm:"not null"`
}
