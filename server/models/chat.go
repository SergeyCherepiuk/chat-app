package models

type Chat struct {
	ID       uint      `json:"id" db:"id" gorm:"primary"`
	Name     string    `json:"name" db:"name"`
	Messages []Message `gorm:"constraint:OnDelete:CASCADE;"`
}
