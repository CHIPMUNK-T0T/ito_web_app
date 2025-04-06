package websock

import (
	"encoding/json"
	"log"
)

// GameMessageHandler は、ゲームのメッセージ処理を行うインターフェース
type GameMessageHandler interface {
	HandleVoteSubmit(userID uint, roomID uint, payload VoteSubmitPayload) error
	// 他のメッセージハンドリングメソッドを追加
}

type MessageHandler struct {
	gameHandler GameMessageHandler
}

func NewMessageHandler(gameHandler GameMessageHandler) *MessageHandler {
	return &MessageHandler{
		gameHandler: gameHandler,
	}
}

func (h *MessageHandler) HandleMessage(client *Client, messageBytes []byte) error {
	var message Message
	if err := json.Unmarshal(messageBytes, &message); err != nil {
		return err
	}

	switch message.Type {
	case MessageTypeVoteSubmit:
		var payload VoteSubmitPayload
		if err := json.Unmarshal(message.Payload, &payload); err != nil {
			return err
		}
		return h.gameHandler.HandleVoteSubmit(client.UserID, client.RoomID, payload)
	default:
		log.Printf("Unknown message type: %s", message.Type)
		return nil
	}
} 