package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/SergeyCherepiuk/chat-app/authentication/models"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type AuthStorage struct {
	pdb *gorm.DB
	rdb *redis.Client
}

func NewAuthStorage(pdb *gorm.DB, rdb *redis.Client) *AuthStorage {
	return &AuthStorage{pdb: pdb, rdb: rdb}
}

func (storage AuthStorage) SignUp(user models.User) (uuid.UUID, error) {
	sessionId := uuid.New()

	tx := storage.pdb.Begin()
	pipe := storage.rdb.Pipeline()

	r := tx.Create(&user)
	if r.Error != nil {
		tx.Rollback()
		pipe.Discard()
		return uuid.UUID{}, r.Error
	}

	err := pipe.Set(
		context.Background(),
		sessionId.String(),
		fmt.Sprint(user.ID),
		7*24*time.Hour,
	).Err()
	if err != nil {
		tx.Rollback()
		pipe.Discard()
		return uuid.UUID{}, err
	}

	tx.Commit()
	pipe.Exec(context.Background())
	return sessionId, nil
}
