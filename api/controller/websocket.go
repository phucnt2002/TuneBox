package controller

import (
	"TuneBox/domain"
	"TuneBox/repository"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strings"
	"sync"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Điều chỉnh trong production
	},
}

type WebSocketController struct {
	groupRepo *repository.GroupRepository
	clients   map[*websocket.Conn]string // Map client đến groupID
	mutex     sync.Mutex
}

func NewWebSocketController(groupRepo *repository.GroupRepository) *WebSocketController {
	return &WebSocketController{
		groupRepo: groupRepo,
		clients:   make(map[*websocket.Conn]string),
	}
}

func (c *WebSocketController) HandleConnections(w http.ResponseWriter, r *http.Request) {
	// Trích xuất groupName từ URL (VD: /party-group -> "party-group")
	path := strings.TrimPrefix(r.URL.Path, "/")
	if path == "" {
		http.Error(w, "Group name is required", http.StatusBadRequest)
		return
	}
	groupName := path

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error upgrading to WebSocket: %v", err)
		return
	}
	defer ws.Close()

	// Tạo hoặc tham gia group dựa trên groupName
	groupID := c.groupRepo.CreateGroup(groupName)

	c.mutex.Lock()
	c.clients[ws] = groupID
	c.mutex.Unlock()

	// Gửi playlist hiện tại của group
	playlist := c.groupRepo.GetGroup(groupID).Playlist
	if err := ws.WriteJSON(map[string]interface{}{"type": "playlist", "data": playlist, "groupId": groupID}); err != nil {
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

		currentGroupID := c.clients[ws]
		switch msg["type"] {
		case "addSong":
			songData := msg["song"].(map[string]interface{})
			song := domain.Song{
				Title:   songData["title"].(string),
				VideoID: songData["videoId"].(string),
			}
			c.groupRepo.AddSong(currentGroupID, song)
			c.broadcastPlaylist(currentGroupID)
		case "removeSong":
			index := int(msg["index"].(float64))
			c.groupRepo.RemoveSong(currentGroupID, index)
			c.broadcastPlaylist(currentGroupID)
		case "playNext":
			song, _ := c.groupRepo.GetNextSong(currentGroupID)
			c.broadcastMessage(currentGroupID, map[string]interface{}{"type": "playSong", "data": song})
			c.broadcastPlaylist(currentGroupID)
		}
	}
}

func (c *WebSocketController) broadcastPlaylist(groupID string) {
	playlist := c.groupRepo.GetGroup(groupID).Playlist
	c.mutex.Lock()
	defer c.mutex.Unlock()
	for client, clientGroupID := range c.clients {
		if clientGroupID == groupID {
			if err := client.WriteJSON(map[string]interface{}{"type": "playlist", "data": playlist, "groupId": groupID}); err != nil {
				log.Printf("Error broadcasting playlist: %v", err)
				client.Close()
				delete(c.clients, client)
			}
		}
	}
}

func (c *WebSocketController) broadcastMessage(groupID string, msg map[string]interface{}) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	for client, clientGroupID := range c.clients {
		if clientGroupID == groupID {
			if err := client.WriteJSON(msg); err != nil {
				log.Printf("Error broadcasting message: %v", err)
				client.Close()
				delete(c.clients, client)
			}
		}
	}
}
