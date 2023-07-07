package auth

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

func (manager SessionManager) CreateSession(userId uint) uuid.UUID {
	sessionId := uuid.New()
	manager.rdb.Set(context.Background(), fmt.Sprint(sessionId), userId, 10*time.Second)
	return sessionId
}

func (manager SessionManager) CheckSession(sessionId uuid.UUID) (uint, error) {
	if manager.rdb.Exists(context.Background(), fmt.Sprint(sessionId)).Val() < 1 {
		return 0, errors.New("session expired")
	}

	userId, err := strconv.ParseUint(manager.rdb.Get(context.Background(), fmt.Sprint(sessionId)).Val(), 10, 64)
	if err != nil {
		return 0, err
	}

	return uint(userId), nil
}

func (manager SessionManager) DeleteSession(sessionId uuid.UUID) {
	manager.rdb.Del(context.Background(), fmt.Sprint(sessionId))
}
