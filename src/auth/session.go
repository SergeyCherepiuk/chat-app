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

func (manager SessionManager) CreateSession(userId uint) (uuid.UUID, error) {
	sessionId := uuid.New()

	oldSessionId := manager.rdb.Get(context.Background(), fmt.Sprint(userId)).Val()

	// Create transaction (pipeline)
	pipe := manager.rdb.TxPipeline()

	// Delete old session
	pipe.Del(context.Background(), oldSessionId)
	pipe.Del(context.Background(), fmt.Sprint(userId))

	// Create new session
	pipe.Set(context.Background(), fmt.Sprint(sessionId), userId, 10*time.Second)
	pipe.Set(context.Background(), fmt.Sprint(userId), fmt.Sprint(sessionId), 10*time.Second)

	// Commit
	pipe.Exec(context.Background())

	return sessionId, nil
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

func (manager SessionManager) DeleteSessions(sessionId uuid.UUID) {
	userId := manager.rdb.Get(context.Background(), fmt.Sprint(sessionId)).Val()
	manager.rdb.Del(context.Background(), userId)
	manager.rdb.Del(context.Background(), fmt.Sprint(sessionId))
}
