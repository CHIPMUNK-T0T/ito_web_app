package usecase

type GameEvent struct {
	Type    GameEventType `json:"type"`
	Payload interface{}   `json:"payload"`
}

type GameEventType string

const (
	EventPlayerJoined   GameEventType = "player_joined"
	EventPlayerLeft     GameEventType = "player_left"
	EventPlayerReady    GameEventType = "player_ready"
	EventGameStarted    GameEventType = "game_started"
	EventThemeAnnounced GameEventType = "theme_announced"
	EventCardDealt      GameEventType = "card_dealt"
	EventVoteStarted    GameEventType = "vote_started"
	EventVoteResult     GameEventType = "vote_result"
	EventCardRevealed   GameEventType = "card_revealed"
	EventGameOver       GameEventType = "game_over"
)

type GameEventHandler interface {
	HandleGameEvent(event GameEvent) error
}