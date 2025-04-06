package websock

import "encoding/json"

type MessageType string

const (
    MessageTypeGameStart   MessageType = "game_start"
    MessageTypeCardDealt   MessageType = "card_dealt"
    MessageTypeVoteStart   MessageType = "vote_start"
    MessageTypeVoteSubmit  MessageType = "vote_submit"
    MessageTypeCardReveal  MessageType = "card_reveal"
    MessageTypeGameEnd     MessageType = "game_end"
)

type Message struct {
    Type    MessageType     `json:"type"`
    RoomID  uint           `json:"room_id"`
    UserID  uint           `json:"user_id"`
    Payload json.RawMessage `json:"payload"`
}

// 各メッセージタイプのペイロード構造体
type CardDealtPayload struct {
    CardNumber int `json:"card_number"`
}

type VoteStartPayload struct {
    TargetUserID uint `json:"target_user_id"`
    TimeoutAt    int64 `json:"timeout_at"` // Unix timestamp
}

type VoteSubmitPayload struct {
    Approve bool `json:"approve"`
}

type CardRevealPayload struct {
    UserID     uint `json:"user_id"`
    CardNumber int  `json:"card_number"`
}

type GameEndPayload struct {
    IsSuccess bool   `json:"is_success"`
    Reason    string `json:"reason,omitempty"`
} 