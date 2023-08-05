package authstorage

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type SessionManager struct {
	rdb *redis.Client
}

func NewSessionManager(rdb *redis.Client) *SessionManager {
	return &SessionManager{rdb: rdb}
}

func (manager SessionManager) Create(userId uint) (uuid.UUID, error) {
	oldSessionId, err := manager.rdb.Get(context.Background(), fmt.Sprint(userId)).Result()
	if err == nil {
		manager.rdb.Del(context.Background(), oldSessionId)
	}

	sessionId := uuid.New()
	pipe := manager.rdb.Pipeline()
	
	pipe.Set(context.Background(), sessionId.String(), fmt.Sprint(userId), 7*24*time.Hour)
	pipe.Set(context.Background(), fmt.Sprint(userId), sessionId.String(), 7*24*time.Hour)
	
	_, err = pipe.Exec(context.Background())
	if err != nil { 
		return uuid.UUID{}, err
	}

	return sessionId, nil
}

func (manager SessionManager) Check(sessionId uuid.UUID) (uint, error) {
	userIdStr, err := manager.rdb.Get(context.Background(), sessionId.String()).Result()
	if err != nil {
		return 0, errors.New("session not found")
	}

	userId, err := strconv.ParseUint(userIdStr, 10, 64)
	if err != nil {
		return 0, err
	}

	return uint(userId), nil
}

func (manager SessionManager) Invalidate(sessionId uuid.UUID) error {
	userId, err := manager.rdb.Get(context.Background(), sessionId.String()).Result()
	if err != nil {
		return err
	}

	pipe := manager.rdb.Pipeline()
	pipe.Del(context.Background(), sessionId.String())
	pipe.Del(context.Background(), userId)
	_, err = pipe.Exec(context.Background())
	return err
}
