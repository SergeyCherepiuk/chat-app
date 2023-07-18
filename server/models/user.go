package models

type User struct {
	ID        uint   `json:"id" db:"id" gorm:"primary"`
	FirstName string `json:"first_name" db:"first_name" gorm:"not null"`
	LastName  string `json:"last_name" db:"last_name" gorm:"not null"`
	Username  string `json:"username" db:"username" gorm:"not null;unique"`
	Password  string `json:"password" db:"password" gorm:"not null"`
	Messages  []Message
}
