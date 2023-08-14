package domain

import (
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/gofiber/contrib/websocket"
)

type ConnectionManagerService[T comparable] interface {
	Connect(key T, conn *websocket.Conn)
	Disconnect(key T, conn *websocket.Conn)
	GetConnections(key T) *hashset.Set
}
