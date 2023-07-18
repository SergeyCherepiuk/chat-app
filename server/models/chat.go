package models

type Chat struct {
	ID       uint `json:"id" db:"id" gorm:"primary"`
	Messages []Message
}
