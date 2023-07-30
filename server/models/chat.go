package models

type Chat struct {
	ID       uint      `json:"id" db:"id" gorm:"primary"`
	Name     string    `json:"name" db:"name" gorm:"not null"`
	Messages []Message `json:"-" gorm:"constraint:OnDelete:CASCADE;"`
}
