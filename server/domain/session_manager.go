package domain

import "github.com/google/uuid"

type SessionManagerService interface {
	Create(userId uint) (uuid.UUID, error)
	Check(sessionId uuid.UUID) (uint, error)
	Invalidate(sessionId uuid.UUID) error
}
