package models

import (
	"time"
)

type User struct {
	ID             uint      `json:"id"`
	FirstName      string    `json:"first_name" db:"first_name"`
	LastName       string    `json:"last_name" db:"last_name"`
	Username       string    `json:"username"`
	Password       string    `json:"password"`
	Description    string    `json:"description"`
	ProfilePicture []byte    `json:"-" db:"profile_picture"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
}
