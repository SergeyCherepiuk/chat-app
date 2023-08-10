package domain

import "github.com/google/uuid"

type AuthService interface {
	SignUp(user User) (uuid.UUID, uint, error)
	Login(username, password string) (uuid.UUID, uint, error)
	Check(sessionId uuid.UUID) (uint, error)
	Logout(sessionId uuid.UUID) error
}
