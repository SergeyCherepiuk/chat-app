package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/SergeyCherepiuk/chat-app/authentication/models"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
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

	err = pipe.Set(
		context.Background(),
		fmt.Sprint(user.ID),
		sessionId.String(),
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

func (storage AuthStorage) Login(username, password string) (uuid.UUID, error) {
	user := models.User{}
	r := storage.pdb.Where("username = ?", username).First(&user)
	if r.Error != nil {
		return uuid.UUID{}, r.Error
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return uuid.UUID{}, err
	}

	oldSessionId, err := storage.rdb.Get(context.Background(), fmt.Sprint(user.ID)).Result()
	if err == nil {
		storage.rdb.Del(context.Background(), oldSessionId)
	}

	sessionId := uuid.New()
	pipe := storage.rdb.Pipeline()

	err = pipe.Set(
		context.Background(),
		sessionId.String(),
		fmt.Sprint(user.ID),
		7*24*time.Hour,
	).Err()
	if err != nil {
		pipe.Discard()
		return uuid.UUID{}, err
	}

	err = pipe.Set(
		context.Background(),
		fmt.Sprint(user.ID),
		sessionId.String(),
		7*24*time.Hour,
	).Err()
	if err != nil {
		pipe.Discard()
		return uuid.UUID{}, err
	}

	pipe.Exec(context.Background())
	return sessionId, nil
}
