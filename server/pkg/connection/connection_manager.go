package connection

import (
	"sync"

	"github.com/emirpasic/gods/sets/hashset"
	"github.com/gofiber/contrib/websocket"
)

type ConnectionManagerService[T comparable] struct {
	mu          sync.RWMutex
	Connections map[T]*hashset.Set
}

func NewConnectionManager[T comparable]() *ConnectionManagerService[T] {
	return &ConnectionManagerService[T]{Connections: make(map[T]*hashset.Set)}
}

func (manager *ConnectionManagerService[T]) Connect(key T, conn *websocket.Conn) {
	manager.mu.Lock()
	defer manager.mu.Unlock()

	if _, ok := manager.Connections[key]; !ok {
		manager.Connections[key] = hashset.New()
	}
	manager.Connections[key].Add(conn)
}

func (manager *ConnectionManagerService[T]) Disconnect(key T, conn *websocket.Conn) {
	manager.mu.Lock()
	defer manager.mu.Unlock()

	if _, ok := manager.Connections[key]; ok {
		manager.Connections[key].Remove(conn)
		if manager.Connections[key].Empty() {
			delete(manager.Connections, key)
		}
	}
}

func (manager *ConnectionManagerService[T]) GetConnections(key T) []*websocket.Conn {
	manager.mu.RLock()
	defer manager.mu.RUnlock()

	connections := make([]*websocket.Conn, manager.Connections[key].Size())
	for i, conn := range manager.Connections[key].Values() {
		connections[i] = conn.(*websocket.Conn)
	}
	return connections
}
