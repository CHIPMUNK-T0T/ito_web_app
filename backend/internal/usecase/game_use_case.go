package usecase

import (
	"encoding/json"
	"errors"
	"math/rand"
	"time"

	domainErrors "CHIPMUNK-T0T/ito_web_app/internal/domain/errors"
	"CHIPMUNK-T0T/ito_web_app/internal/entity/domain"
	"CHIPMUNK-T0T/ito_web_app/internal/websock"

	// "CHIPMUNK-T0T/ito_web_app/internal/websock"

	gorillaws "github.com/gorilla/websocket"
)

type GameUseCase struct {
	roomUseCase  *RoomUseCase
	userUseCase  *UserUseCase
	themeUseCase *ThemeUseCase
	gameHub      *websock.GameHub
}

func NewGameUseCase(roomUseCase *RoomUseCase, userUseCase *UserUseCase, themeUseCase *ThemeUseCase, gameHub *websock.GameHub) *GameUseCase {
	return &GameUseCase{
		roomUseCase:  roomUseCase,
		userUseCase:  userUseCase,
		themeUseCase: themeUseCase,
		gameHub:      gameHub,
	}
}

func (uc *GameUseCase) PrepareGame(roomID uint) (*domain.Game, error) {
	room, err := uc.roomUseCase.GetRoom(roomID)
	if err != nil {
		return nil, err
	}

	if room.GetPlayerCount() < 2 {
		return nil, errors.New("ゲームを開始するには2人以上のプレイヤーが必要です")
	}

	players := room.GetPlayers()
	for _, player := range players {
		if !player.IsReady() {
			return nil, errors.New("全員の準備が完了していません")
		}
	}

	return domain.NewGame(*room), nil
}

func (uc *GameUseCase) StartRound(game *domain.Game) error {
	if game.Status() != domain.GameStatusWaiting {
		return errors.New("ゲームの状態が不正です")
	}

	theme, err := uc.themeUseCase.GetRandomTheme()
	if err != nil {
		return err
	}

	// カードを配布
	room := game.Room()
	players := room.GetPlayers()
	numbers := make([]int, len(players))
	for i := range numbers {
		numbers[i] = i + 1
	}

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(numbers), func(i, j int) {
		numbers[i], numbers[j] = numbers[j], numbers[i]
	})

	for i, player := range players {
		player.SetCardValue(numbers[i])
	}

	game.SetTheme(theme)
	game.StartPlaying()

	// プレイヤーにカードを通知
	for _, player := range players {
		payload := websock.CardDealtPayload{
			CardNumber: *player.CardValue(),
		}
		user := player.User()
		payloadBytes, err := json.Marshal(payload)
		if err != nil {
			return err
		}
		msg, err := json.Marshal(websock.Message{
			Type:    websock.MessageTypeCardDealt,
			RoomID:  room.ID(),
			UserID:  user.ID(),
			Payload: json.RawMessage(payloadBytes),
		})
		if err != nil {
			return err
		}
		uc.gameHub.SendToUser(user.ID(), msg)
	}

	return nil
}

// 数字公開の投票を開始
func (uc *GameUseCase) InitiateVote(game *domain.Game, userID uint) error {
	if game.Status() != domain.GameStatusPlaying {
		return errors.New("ゲームの状態が不正です")
	}

	room := game.Room()
	players := room.GetPlayers()
	var targetPlayer *domain.RoomPlayer

	for _, player := range players {
		user := player.User()
		if user.ID() == userID {
			targetPlayer = &player
			break
		}
	}

	if targetPlayer == nil {
		return errors.New("プレイヤーが見つかりません")
	}

	if targetPlayer.IsReady() {
		return errors.New("すでにカードは公開済みです")
	}

	return game.InitiateVote(userID)
}

// 投票を処理
func (uc *GameUseCase) ProcessVote(game *domain.Game, voterID uint, approve bool) error {
	if game.VoteState() == nil {
		return errors.New("投票が開始されていません")
	}

	if game.IsVoteTimeout() {
		return uc.finalizeVote(game)
	}

	return game.ProcessVote(voterID, approve)
}

// 投票を確定
func (uc *GameUseCase) finalizeVote(game *domain.Game) error {
	if !game.IsVoteSuccessful() {
		game.ClearVoteState()
		return errors.New("投票が却下されました")
	}

	return uc.revealCard(game, game.VoteState().TargetPlayerID())
}

// カードを公開
func (uc *GameUseCase) revealCard(game *domain.Game, userID uint) error {
	room := game.Room()
	players := room.GetPlayers()
	var currentPlayer *domain.RoomPlayer

	for _, player := range players {
		user := player.User()
		if user.ID() == userID {
			currentPlayer = &player
			break
		}
	}

	if currentPlayer == nil {
		return errors.New("プレイヤーが見つかりません")
	}

	// より小さい数字を持つプレイヤーがいないか確認
	for _, player := range players {
		if !player.IsReady() && *player.CardValue() < *currentPlayer.CardValue() {
			game.Finish()
			return errors.New("より小さい数字のカードが残っています")
		}
	}

	currentPlayer.SetIsReady(true)
	game.ClearVoteState()

	// 全員のカードが公開されたらゲーム終了
	allRevealed := true
	for _, player := range players {
		if !player.IsReady() {
			allRevealed = false
			break
		}
	}
	if allRevealed {
		game.Finish()
	}

	return nil
}

func (uc *GameUseCase) HandlePlayerConnection(userID, roomID uint, conn *gorillaws.Conn) error {
	client := &websock.Client{
		Hub:    uc.gameHub,
		Conn:   conn,
		Send:   make(chan []byte, 256),
		UserID: userID,
		RoomID: roomID,
	}

	uc.gameHub.Register <- client

	go client.ReadPump()
	go client.WritePump()

	return nil
}

func (uc *GameUseCase) HandlePlayerDisconnect(game *domain.Game, userID uint) error {
	room := game.Room()
	players := room.GetPlayers()
	for _, player := range players {
		user := player.User()
		if user.ID() == userID {
			player.SetIsConnected(false)
			return nil
		}
	}
	return domainErrors.ErrPlayerNotFound
}

// GameUseCaseにwebsock.GameMessageHandlerインターフェースを実装
func (uc *GameUseCase) HandleVoteSubmit(userID uint, roomID uint, payload websock.VoteSubmitPayload) error {
	// ゲームの取得
	game, err := uc.GetGameByRoomID(roomID)
	if err != nil {
		return err
	}

	// 投票の処理
	return uc.ProcessVote(game, userID, payload.Approve)
}

// GetGameByRoomIDを追加
func (uc *GameUseCase) GetGameByRoomID(roomID uint) (*domain.Game, error) {
	room, err := uc.roomUseCase.GetRoom(roomID)
	if err != nil {
		return nil, err
	}

	return domain.NewGame(*room), nil
}

func (uc *GameUseCase) SetPlayerReady(roomID, userID uint) error {
	game, err := uc.GetGameByRoomID(roomID)
	if err != nil {
		return err
	}

	return game.SetPlayerReady(userID)
}

func (uc *GameUseCase) StartGame(roomID uint) error {
	game, err := uc.GetGameByRoomID(roomID)
	if err != nil {
		return err
	}

	room := game.Room()
	if room.GetPlayerCount() < 2 {
		return errors.New("ゲームを開始するには2人以上のプレイヤーが必要です")
	}

	// カードの配布
	if err := game.DealCards(); err != nil {
		return err
	}

	// お題の設定
	theme, err := uc.themeUseCase.GetRandomTheme()
	if err != nil {
		return err
	}
	game.SetTheme(theme)

	// ゲーム開始状態に変更
	game.Start()

	return nil
} 