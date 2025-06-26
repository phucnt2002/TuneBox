package controller

import (
	"TuneBox/repository"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WebSocketController struct {
	repo    *repository.InMemoryRepository
	clients map[*websocket.Conn]bool
	mutex   sync.Mutex
}

func NewWebSocketController(repo *repository.InMemoryRepository) *WebSocketController {
	return &WebSocketController{
		repo:    repo,
		clients: make(map[*websocket.Conn]bool),
	}
}

func (c *WebSocketController) HandleConnection(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading to WebSocket: %v", err)
		return
	}
	defer ws.Close()

	c.mutex.Lock()
	c.clients[ws] = true
	c.mutex.Unlock()

}
