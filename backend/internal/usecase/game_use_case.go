package usecase

import (
	"encoding/json"
	"errors"
	"math/rand"
	"sync"
	"time"

	"CHIPMUNK-T0T/ito_web_app/internal/entity/domain"
	"CHIPMUNK-T0T/ito_web_app/internal/websock"
	"github.com/gorilla/websocket"
)

// InMemoryHand はプレイヤー1人の手札情報をメモリ上で保持
type InMemoryHand struct {
	UserID     uint
	CardValue  int
	IsRevealed bool
}

// InMemorySession は特定のルームのゲーム進行状況をメモリ上で保持
type InMemorySession struct {
	RoomID       uint
	Status       string // "waiting", "playing", "finished"
	ThemeContent string
	Hands        map[uint]*InMemoryHand // UserID -> Hand
	CardPool     []int                  // 現在の山札
	UsedCards    []int                  // 使用済み（トラッシュ）
	mu           sync.RWMutex
}

type GameUseCase struct {
	roomUseCase  *RoomUseCase
	userUseCase  *UserUseCase
	themeUseCase *ThemeUseCase
	gameHub      *websock.GameHub
	
	// メモリ上でのゲームセッション管理
	sessions   map[uint]*InMemorySession
	sessionsMu sync.RWMutex
}

func NewGameUseCase(roomUseCase *RoomUseCase, userUseCase *UserUseCase, themeUseCase *ThemeUseCase, gameHub *websock.GameHub) *GameUseCase {
	return &GameUseCase{
		roomUseCase:  roomUseCase,
		userUseCase:  userUseCase,
		themeUseCase: themeUseCase,
		gameHub:      gameHub,
		sessions:     make(map[uint]*InMemorySession),
	}
}

// getOrCreateSession はルームIDに対応するメモリ上のセッションを取得、なければ作成
func (uc *GameUseCase) getOrCreateSession(roomID uint) *InMemorySession {
	uc.sessionsMu.Lock()
	defer uc.sessionsMu.Unlock()

	if session, ok := uc.sessions[roomID]; ok {
		return session
	}

	// 1-100の初期山札を作成
	pool := make([]int, 100)
	for i := 0; i < 100; i++ {
		pool[i] = i + 1
	}

	session := &InMemorySession{
		RoomID:    roomID,
		Status:    "waiting",
		Hands:     make(map[uint]*InMemoryHand),
		CardPool:  pool,
		UsedCards: make([]int, 0),
	}
	uc.sessions[roomID] = session
	return session
}

func (uc *GameUseCase) SetPlayerReady(roomID, userID uint) error {
	room, err := uc.roomUseCase.GetRoom(roomID)
	if err != nil {
		return err
	}

	session := uc.getOrCreateSession(roomID)
	session.mu.Lock()
	defer session.mu.Unlock()

	// 部屋のプレイヤー情報を更新
	players := room.GetPlayers()
	found := false
	for i := range players {
		if players[i].User().ID() == userID {
			players[i].SetIsReady(true)
			found = true
			break
		}
	}
	if !found {
		return errors.New("プレイヤーが見つかりません")
	}

	// 永続化（メモリ上のリポジトリに反映）
	if err := uc.roomUseCase.roomRepo.Update(room); err != nil {
		return err
	}

	// 全員に通知
	readyPayload := websock.PlayerReadyPayload{
		UserID:  userID,
		IsReady: true,
	}
	payloadBytes, _ := json.Marshal(readyPayload)
	msg, _ := json.Marshal(websock.Message{
		Type:    websock.MessageTypePlayerReady,
		RoomID:  roomID,
		UserID:  userID,
		Payload: json.RawMessage(payloadBytes),
	})
	uc.gameHub.BroadcastToRoom(roomID, msg)

	return nil
}
func (uc *GameUseCase) StartGame(roomID uint, customTheme string) error {
	room, err := uc.roomUseCase.GetRoom(roomID)
	if err != nil {
		return err
	}

	players := room.GetPlayers()
	playerCount := len(players)
	if playerCount < 2 {
		return errors.New("ゲームを開始するには2人以上のプレイヤーが必要です")
	}

	// 全員の準備完了をチェック
	for _, p := range players {
		if !p.IsReady() {
			return errors.New("全員が準備完了になるまで開始できません")
		}
	}

	session := uc.getOrCreateSession(roomID)
	session.mu.Lock()
	defer session.mu.Unlock()

	// 自動リフレッシュ: 山札が足りない場合はリセット
	if len(session.CardPool) < playerCount {
		pool := make([]int, 100)
		for i := 0; i < 100; i++ {
			pool[i] = i + 1
		}
		session.CardPool = pool
		session.UsedCards = make([]int, 0)
	}

	// お題の設定
	theme := customTheme
	if theme == "" {
		theme, err = uc.themeUseCase.GetRandomTheme()
		if err != nil {
			return err
		}
	}
	session.ThemeContent = theme
	session.Status = "playing"

	// カードの配布 (現在のプールをシャッフルして抽出)
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)
	
	r.Shuffle(len(session.CardPool), func(i, j int) {
		session.CardPool[i], session.CardPool[j] = session.CardPool[j], session.CardPool[i]
	})

	// 人数分抽出
	dealtCards := session.CardPool[:playerCount]
	session.CardPool = session.CardPool[playerCount:]
	session.UsedCards = append(session.UsedCards, dealtCards...)

	players = room.GetPlayers()
	session.Hands = make(map[uint]*InMemoryHand)
	for i, player := range players {
		uID := player.User().ID()
		cardVal := dealtCards[i]
		session.Hands[uID] = &InMemoryHand{
			UserID:     uID,
			CardValue:  cardVal,
			IsRevealed: false,
		}

		// 個別通知
		cardPayload := websock.CardDealtPayload{CardNumber: cardVal}
		cardBytes, _ := json.Marshal(cardPayload)
		cardMsg, _ := json.Marshal(websock.Message{
			Type:    websock.MessageTypeCardDealt,
			RoomID:  roomID,
			UserID:  uID,
			Payload: json.RawMessage(cardBytes),
		})
		uc.gameHub.SendToUser(uID, cardMsg)
	}

	// 全員にゲーム開始を通知
	startPayload := websock.GameStartPayload{Theme: theme}
	payloadBytes, _ := json.Marshal(startPayload)
	msg, _ := json.Marshal(websock.Message{
		Type:    websock.MessageTypeGameStart,
		RoomID:  roomID,
		Payload: json.RawMessage(payloadBytes),
	})
	uc.gameHub.BroadcastToRoom(roomID, msg)

	// 全員に山札情報を通知
	deckPayload := map[string]interface{}{"remaining_count": len(session.CardPool), "total_count": 100}
	deckBytes, _ := json.Marshal(deckPayload)
	deckMsg, _ := json.Marshal(websock.Message{
		Type:    "deck_info",
		RoomID:  roomID,
		Payload: json.RawMessage(deckBytes),
	})
	uc.gameHub.BroadcastToRoom(roomID, deckMsg)

	return nil
}

// HandlePlayCard は websock.GameMessageHandler インターフェースを実装
func (uc *GameUseCase) HandlePlayCard(userID uint, roomID uint, payload websock.PlayCardPayload) error {
	session := uc.getOrCreateSession(roomID)
	session.mu.Lock()
	defer session.mu.Unlock()

	if session.Status != "playing" {
		return errors.New("ゲーム中ではありません")
	}

	hand, ok := session.Hands[userID]
	if !ok {
		return errors.New("プレイヤーの手札が見つかりません")
	}

	if hand.IsRevealed {
		return errors.New("すでにカードは提出済みです")
	}

	// 判定ロジック
	for _, h := range session.Hands {
		if !h.IsRevealed && h.UserID != userID && h.CardValue < hand.CardValue {
			session.Status = "finished"
			
			resultPayload := websock.GameResultPayload{
				Status:  "failed",
				Message: "失敗！より小さい数字を持っているプレイヤーがいました",
			}
			payloadBytes, _ := json.Marshal(resultPayload)
			msg, _ := json.Marshal(websock.Message{
				Type:    websock.MessageTypeGameResult,
				RoomID:  roomID,
				Payload: json.RawMessage(payloadBytes),
			})
			uc.gameHub.BroadcastToRoom(roomID, msg)
			return nil
		}
	}

	hand.IsRevealed = true

	playPayload := websock.PlayCardPayload{
		UserID:     userID,
		CardNumber: hand.CardValue,
	}
	payloadBytes, _ := json.Marshal(playPayload)
	msg, _ := json.Marshal(websock.Message{
		Type:    websock.MessageTypePlayCard,
		RoomID:  roomID,
		Payload: json.RawMessage(payloadBytes),
	})
	uc.gameHub.BroadcastToRoom(roomID, msg)

	allRevealed := true
	for _, h := range session.Hands {
		if !h.IsRevealed {
			allRevealed = false
			break
		}
	}

	if allRevealed {
		session.Status = "finished"
		resultPayload := websock.GameResultPayload{
			Status:  "success",
			Message: "完全成功！全員のカードが正しい順番で出されました",
		}
		resultBytes, _ := json.Marshal(resultPayload)
		successMsg, _ := json.Marshal(websock.Message{
			Type:    websock.MessageTypeGameResult,
			RoomID:  roomID,
			Payload: json.RawMessage(resultBytes),
		})
		uc.gameHub.BroadcastToRoom(roomID, successMsg)
	}

	return nil
}

func (uc *GameUseCase) GetGameStatus(roomID uint) (map[string]interface{}, error) {
	session := uc.getOrCreateSession(roomID)
	session.mu.RLock()
	defer session.mu.RUnlock()

	return map[string]interface{}{
		"status": session.Status,
		"theme":  session.ThemeContent,
	}, nil
}

func (uc *GameUseCase) HandlePlayerConnection(userID, roomID uint, conn *websocket.Conn) error {
	client := websock.NewClient(uc.gameHub, conn, userID, roomID)
	uc.gameHub.Register <- client

	// プレイヤーが接続したことを通知
	uc.updatePlayerConnectionStatus(userID, roomID, true)

	go client.WritePump()
	
	// ReadPump を別の goroutine で実行し、切断時にステータスを更新する
	go func() {
		defer conn.Close()
		client.ReadPump()
		// 切断時の処理
		uc.updatePlayerConnectionStatus(userID, roomID, false)
	}()

	return nil
}

func (uc *GameUseCase) updatePlayerConnectionStatus(userID, roomID uint, isConnected bool) {
	room, err := uc.roomUseCase.GetRoom(roomID)
	if err != nil {
		return
	}

	session := uc.getOrCreateSession(roomID)
	session.mu.Lock()
	defer session.mu.Unlock()

	if isConnected {
		// 接続時のステータス更新
		players := room.GetPlayers()
		for i := range players {
			if players[i].User().ID() == userID {
				players[i].SetIsConnected(true)
				break
			}
		}
		uc.roomUseCase.roomRepo.Update(room)
	} else {
		// 切断時はルームとセッションから完全に削除
		room.RemovePlayer(userID)
		uc.roomUseCase.roomRepo.Update(room)
		
		// ゲームセッション（手札）からも削除
		delete(session.Hands, userID)

		// もしゲーム中（playing）であれば、残りの人数をチェック
		if session.Status == "playing" {
			// プレイヤーが2人未満になったら強制終了
			if len(session.Hands) < 2 {
				session.Status = "finished"
				resultPayload := websock.GameResultPayload{
					Status:  "failed",
					Message: "プレイヤーが2人未満になったため、ゲームを終了しました",
				}
				payloadBytes, _ := json.Marshal(resultPayload)
				msg, _ := json.Marshal(websock.Message{
					Type:    websock.MessageTypeGameResult,
					RoomID:  roomID,
					Payload: json.RawMessage(payloadBytes),
				})
				uc.gameHub.BroadcastToRoom(roomID, msg)
			} else {
				// まだ2人以上いれば、残りのメンバーでクリアしていないか判定
				allRevealed := true
				for _, h := range session.Hands {
					if !h.IsRevealed {
						allRevealed = false
						break
					}
				}

				if allRevealed {
					session.Status = "finished"
					resultPayload := websock.GameResultPayload{
						Status:  "success",
						Message: "完全成功！全員のカードが正しい順番で出されました（離脱者がいたため調整されました）",
					}
					resultBytes, _ := json.Marshal(resultPayload)
					successMsg, _ := json.Marshal(websock.Message{
						Type:    websock.MessageTypeGameResult,
						RoomID:  roomID,
						Payload: json.RawMessage(resultBytes),
					})
					uc.gameHub.BroadcastToRoom(roomID, successMsg)
				}
			}
		}
	}

	// ルーム情報の更新を全員に通知（プレイヤーリストの同期）
	// room_handler で定義されている newRoomResponse と同等の情報を送る必要があるが
	// ここではシンプルにプレイヤー接続/切断イベントまたはルーム更新イベントを送る
	// クライアント側は fetchRoomInfo を定期的に呼んでいるため、WebSocket通知があれば即時反映される
	
	connPayload := websock.PlayerConnectionPayload{
		UserID:      userID,
		IsConnected: isConnected,
	}
	payloadBytes, _ := json.Marshal(connPayload)
	msg, _ := json.Marshal(websock.Message{
		Type:    websock.MessageTypePlayerConnection,
		RoomID:  roomID,
		UserID:  userID,
		Payload: json.RawMessage(payloadBytes),
	})
	uc.gameHub.BroadcastToRoom(roomID, msg)
}

func (uc *GameUseCase) HandleChatMessage(userID uint, roomID uint, payload websock.ChatPayload) error {
	user, err := uc.userUseCase.GetUserByID(userID)
	if err != nil {
		return err
	}

	chatPayload := websock.ChatPayload{
		Username: user.Username(),
		Message:  payload.Message,
		SentAt:   time.Now().Unix(),
	}
	payloadBytes, _ := json.Marshal(chatPayload)
	msg, _ := json.Marshal(websock.Message{
		Type:    websock.MessageTypeChat,
		RoomID:  roomID,
		Payload: json.RawMessage(payloadBytes),
	})
	uc.gameHub.BroadcastToRoom(roomID, msg)
	return nil
}

func (uc *GameUseCase) GetGameByRoomID(roomID uint) (*domain.Game, error) {
	room, err := uc.roomUseCase.GetRoom(roomID)
	if err != nil {
		return nil, err
	}

	session := uc.getOrCreateSession(roomID)
	session.mu.RLock()
	defer session.mu.RUnlock()

	game := domain.NewGame(*room)
	game.SetTheme(session.ThemeContent)
	
	if session.Status == "playing" {
		game.StartPlaying()
	} else if session.Status == "finished" {
		game.Finish()
	}

	for _, player := range room.GetPlayers() {
		if hand, ok := session.Hands[player.User().ID()]; ok {
			player.SetCardValue(hand.CardValue)
			player.SetIsReady(hand.IsRevealed)
		}
	}

	return game, nil
}

func (uc *GameUseCase) RefreshDeck(roomID uint) error {
	session := uc.getOrCreateSession(roomID)
	session.mu.Lock()
	defer session.mu.Unlock()

	// 1-100の山札をリセット
	pool := make([]int, 100)
	for i := 0; i < 100; i++ {
		pool[i] = i + 1
	}
	session.CardPool = pool
	session.UsedCards = make([]int, 0)

	// 全員に山札情報を通知
	deckPayload := map[string]interface{}{"remaining_count": len(session.CardPool), "total_count": 100}
	deckBytes, _ := json.Marshal(deckPayload)
	deckMsg, _ := json.Marshal(websock.Message{
		Type:    "deck_info",
		RoomID:  roomID,
		Payload: json.RawMessage(deckBytes),
	})
	uc.gameHub.BroadcastToRoom(roomID, deckMsg)

	return nil
}
