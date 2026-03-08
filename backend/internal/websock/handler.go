package websock

import (
	"encoding/json"
	"log"
)

// GameMessageHandler は、ゲームのメッセージ処理を行うインターフェース
type GameMessageHandler interface {
	HandleChatMessage(userID uint, roomID uint, payload ChatPayload) error
	HandlePlayCard(userID uint, roomID uint, payload PlayCardPayload) error
}

type MessageHandler struct {
	gameHandler GameMessageHandler
}

func NewMessageHandler(gameHandler GameMessageHandler) *MessageHandler {
	return &MessageHandler{
		gameHandler: gameHandler,
	}
}

func (h *MessageHandler) HandlePlayCard(userID uint, roomID uint, payload PlayCardPayload) error {
	return h.gameHandler.HandlePlayCard(userID, roomID, payload)
}

func (h *MessageHandler) HandleMessage(client *Client, messageBytes []byte) error {
	var message Message
	if err := json.Unmarshal(messageBytes, &message); err != nil {
		return err
	}

	switch message.Type {
	case MessageTypePlayCard:
		var payload PlayCardPayload
		if err := json.Unmarshal(message.Payload, &payload); err != nil {
			return err
		}
		return h.gameHandler.HandlePlayCard(client.UserID, client.RoomID, payload)
	case MessageTypeChat:
		var payload ChatPayload
		if err := json.Unmarshal(message.Payload, &payload); err != nil {
			return err
		}
		return h.gameHandler.HandleChatMessage(client.UserID, client.RoomID, payload)

	default:
		log.Printf("Unknown message type: %s", message.Type)
		return nil
	}
}
