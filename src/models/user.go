package models

type User struct {
	ID       uint      `db:"id"`
	Username string    `json:"username" db:"username" gorm:"unique"`
	Password string    `json:"password" db:"password"`
	Sessions []Session `db:"sessions" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
