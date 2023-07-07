package models

type User struct {
	ID       uint    `json:"id" db:"id"`
	Username string  `json:"username" db:"username" gorm:"unique"`
	Password string  `json:"password" db:"password"`
	Roles    []*Role `gorm:"many2many:user_roles"`
}
