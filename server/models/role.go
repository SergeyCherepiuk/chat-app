package models

type Role struct {
	ID    uint    `json:"id" db:"id"`
	Name  string  `json:"name" db:"name" gorm:"unique"`
	Users []*User `json:"-" gorm:"many2many:user_roles"`
}
