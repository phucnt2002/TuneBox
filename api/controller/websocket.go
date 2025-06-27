package controller

import (
	"TuneBox/domain"
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

func (c *WebSocketController) HandleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error upgrading to WebSocket: %v", err)
		return
	}
	defer ws.Close()

	c.mutex.Lock()
	c.clients[ws] = true
	c.mutex.Unlock()

	// Gửi playlist hiện tại cho client mới
	playlist := c.repo.GetPlayList() // Sửa lỗi chính tả nếu cần: GetPlayList -> GetPlaylist
	if err := ws.WriteJSON(map[string]interface{}{"type": "playlist", "data": playlist}); err != nil {
		log.Printf("Error sending initial playlist: %v", err)
	}

	for {
		var msg map[string]interface{}
		if err := ws.ReadJSON(&msg); err != nil {
			log.Printf("Error reading message: %v", err)
			c.mutex.Lock()
			delete(c.clients, ws)
			c.mutex.Unlock()
			break
		}

		switch msg["type"] {
		case "addSong":
			songData := msg["song"].(map[string]interface{})
			song := domain.Song{
				Title:   songData["title"].(string),
				VideoId: songData["videoId"].(string),
			}
			c.repo.AddSong(song)
			c.broadcastPlaylist()
		case "removeSong":
			index := int(msg["index"].(float64))
			c.repo.RemoveSong(index)
			c.broadcastPlaylist()
		case "playNext":
			song, _ := c.repo.GetNextSong()
			c.broadcastMessage(map[string]interface{}{"type": "playSong", "data": song})
			c.broadcastPlaylist() // Cập nhật danh sách phát mới
		}
	}
}
func (c *WebSocketController) broadcastPlaylist() {
	playlist := c.repo.GetPlayList()
	c.mutex.Lock()
	defer c.mutex.Unlock()
	for client := range c.clients {
		if err := client.WriteJSON(map[string]interface{}{"type": "playlist", "data": playlist}); err != nil {
			log.Printf("Error broadcasting playlist: %v", err)
			client.Close()
			delete(c.clients, client)
		}
	}
}

func (c *WebSocketController) broadcastMessage(msg map[string]interface{}) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	for client := range c.clients {
		if err := client.WriteJSON(msg); err != nil {
			log.Printf("Error broadcasting message: %v", err)
			client.Close()
			delete(c.clients, client)
		}
	}
}
