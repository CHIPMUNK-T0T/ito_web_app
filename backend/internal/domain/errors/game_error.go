package errors

type GameError struct {
	Code    string
	Message string
}

func (e GameError) Error() string {
	return e.Message
}

var (
	ErrGameNotFound       = GameError{Code: "GAME_NOT_FOUND", Message: "ゲームが見つかりません"}
	ErrInvalidGameState   = GameError{Code: "INVALID_GAME_STATE", Message: "ゲームの状態が不正です"}
	ErrPlayerNotFound     = GameError{Code: "PLAYER_NOT_FOUND", Message: "プレイヤーが見つかりません"}
	ErrVoteInProgress     = GameError{Code: "VOTE_IN_PROGRESS", Message: "投票が進行中です"}
	ErrNotEnoughPlayers   = GameError{Code: "NOT_ENOUGH_PLAYERS", Message: "プレイヤー数が不足しています"}
	ErrPlayersNotReady    = GameError{Code: "PLAYERS_NOT_READY", Message: "全員の準備が完了していません"}
	ErrPlayerDisconnected = GameError{Code: "PLAYER_DISCONNECTED", Message: "プレイヤーが切断されています"}
)