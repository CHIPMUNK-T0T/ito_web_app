package websock

import (
	"CHIPMUNK-T0T/ito_web_app/internal/entity/domain"
	"sync"
)

type GameHub struct {
	Clients        map[*Client]bool
	Register       chan *Client
	Unregister     chan *Client
	Broadcast      chan []byte
	Games          map[uint]*domain.Game
	mu             sync.RWMutex
	messageHandler *MessageHandler
}

func NewGameHub() *GameHub {
	return &GameHub{
		Clients:    make(map[*Client]bool),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan []byte),
		Games:      make(map[uint]*domain.Game),
	}
}

func (h *GameHub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.mu.Lock()
			h.Clients[client] = true
			h.mu.Unlock()
		case client := <-h.Unregister:
			h.mu.Lock()
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.Send)
			}
			h.mu.Unlock()
		case message := <-h.Broadcast:
			h.mu.RLock()
			for client := range h.Clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.Clients, client)
				}
			}
			h.mu.RUnlock()
		}
	}
}

func (h *GameHub) BroadcastToRoom(roomID uint, message []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for client := range h.Clients {
		if client.RoomID == roomID {
			select {
			case client.Send <- message:
			default:
				close(client.Send)
				delete(h.Clients, client)
			}
		}
	}
}

func (h *GameHub) SendToUser(userID uint, message []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for client := range h.Clients {
		if client.UserID == userID {
			select {
			case client.Send <- message:
			default:
				close(client.Send)
				delete(h.Clients, client)
			}
		}
	}
}

func (h *GameHub) RegisterMessageHandler(handler *MessageHandler) {
	h.messageHandler = handler
}

func (h *GameHub) MessageHandler() *MessageHandler {
	return h.messageHandler
}
