package websock

import "encoding/json"

type MessageType string

const (
	MessageTypeGameStart         MessageType = "game_start"
	MessageTypeCardDealt         MessageType = "card_dealt"
	MessageTypePlayCard          MessageType = "play_card"
	MessageTypeGameResult        MessageType = "game_result"
	MessageTypeChat              MessageType = "chat_message"
	MessageTypePlayerReady       MessageType = "player_ready"
	MessageTypePlayerConnection  MessageType = "player_connection"
)

type Message struct {
	Type    MessageType     `json:"type"`
	RoomID  uint            `json:"room_id"`
	UserID  uint            `json:"user_id"`
	Payload json.RawMessage `json:"payload"`
}

// 各メッセージタイプのペイロード構造体
type CardDealtPayload struct {
	CardNumber int `json:"card_number"`
}

type PlayCardPayload struct {
	UserID     uint `json:"user_id"`
	CardNumber int  `json:"card_number"`
}

type GameResultPayload struct {
	Status  string `json:"status"` // "success" or "failed"
	Message string `json:"message"`
}

type GameStartPayload struct {
	Theme string `json:"theme"`
}

type ChatPayload struct {
	Username string `json:"username"`
	Message  string `json:"message"`
	SentAt   int64  `json:"sent_at"`
}

type PlayerReadyPayload struct {
	UserID  uint `json:"user_id"`
	IsReady bool `json:"is_ready"`
}

type PlayerConnectionPayload struct {
	UserID      uint `json:"user_id"`
	IsConnected bool `json:"is_connected"`
}


