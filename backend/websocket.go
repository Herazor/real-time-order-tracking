package main

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var (
	clients   = make(map[uuid.UUID]map[*websocket.Conn]bool)
	clientsMu sync.Mutex
)

func HandleOrderWebSocket(c *gin.Context) {
	orderID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID pesanan tidak valid"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	clientsMu.Lock()

	if clients[orderID] == nil {
		clients[orderID] = make(map[*websocket.Conn]bool)
	}

	clients[orderID][conn] = true
	clientsMu.Unlock()

	defer func() {
		clientsMu.Lock()
		delete(clients[orderID], conn)

		if len(clients[orderID]) == 0 {
			delete(clients, orderID)
		}

		clientsMu.Unlock()
		conn.Close()
	}()

	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			break
		}
	}
}

func BroadcastOrderEvent(orderID uuid.UUID, event OrderEvent) {
	clientsMu.Lock()
	defer clientsMu.Unlock()

	for conn := range clients[orderID] {
		if err := conn.WriteJSON(event); err != nil {
			conn.Close()
			delete(clients[orderID], conn)
		}
	}
}