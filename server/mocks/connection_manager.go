package mocks

import (
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/gofiber/contrib/websocket"
)

type ConnectionManagerService[T comparable] struct{}

func NewConnectionManagerService[T comparable]() *ConnectionManagerService[T] {
	return &ConnectionManagerService[T]{}
}

func (service ConnectionManagerService[T]) Connect(key T, conn *websocket.Conn) {}

func (service ConnectionManagerService[T]) Disconnect(key T, conn *websocket.Conn) {}

func (service ConnectionManagerService[T]) GetConnections(key T) *hashset.Set {
	return hashset.New()
}