package router

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// 開発環境では全てのオリジンを許可
		// 本番環境では適切に制限する必要がある
		return true
	},
}

func (r *Router) SetPlayerReady(c *gin.Context) {
	roomID, err := strconv.ParseUint(c.Param("roomId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無効なルームIDです"})
		return
	}

	userID := c.GetUint("user_id")
	if err := r.gameUseCase.SetPlayerReady(uint(roomID), userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "準備完了"})
}

func (r *Router) StartGame(c *gin.Context) {
	roomID, err := strconv.ParseUint(c.Param("roomId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無効なルームIDです"})
		return
	}

	// ルームの存在確認とプレイヤー数チェック
	room, err := r.roomUseCase.GetRoom(uint(roomID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ルームが見つかりません"})
		return
	}

	// ゲーム開始の権限チェック
	userID := c.GetUint("user_id")
	if room.GetCreatorID() != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "ゲームを開始する権限がありません"})
		return
	}

	// ゲーム開始
	if err := r.gameUseCase.StartGame(uint(roomID)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ゲームを開始しました"})
}

func (r *Router) InitiateVote(c *gin.Context) {
	roomID, err := strconv.ParseUint(c.Param("roomId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無効なルームID"})
		return
	}

	var req struct {
		UserID uint `json:"user_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無効なリクエスト"})
		return
	}

	game, err := r.gameUseCase.GetGameByRoomID(uint(roomID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ゲームが見つかりません"})
		return
	}

	if err := r.gameUseCase.InitiateVote(game, req.UserID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "投票を開始しました"})
}

func (r *Router) GetGameStatus(c *gin.Context) {
	roomID, err := strconv.ParseUint(c.Param("roomId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無効なルームID"})
		return
	}

	game, err := r.gameUseCase.GetGameByRoomID(uint(roomID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ゲームが見つかりません"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": game.Status(),
		"theme":  game.Theme(),
		"room":   game.Room(),
	})
}

func (r *Router) HandleWebSocket(c *gin.Context) {
	// ユーザー認証の確認
	userID := c.GetUint("user_id")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "認証が必要です"})
		return
	}

	// ルームIDの取得と検証
	roomID, err := strconv.ParseUint(c.Param("roomId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無効なルームIDです"})
		return
	}

	// ルームの存在確認
	room, err := r.roomUseCase.GetRoom(uint(roomID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ルームが見つかりません"})
		return
	}

	// ユーザーがルームに所属しているか確認
	if !room.HasPlayer(userID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "このルームに参加していません"})
		return
	}

	// WebSocket接続へのアップグレード
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "WebSocket接続に失敗しました"})
		return
	}
	defer conn.Close()

	// WebSocket接続の確立を通知
	if err := conn.WriteJSON(gin.H{
		"type": "connected",
		"data": gin.H{
			"user_id": userID,
			"room_id": roomID,
		},
	}); err != nil {
		return
	}

	// ゲームロジックの処理
	if err := r.gameUseCase.HandlePlayerConnection(userID, uint(roomID), conn); err != nil {
		conn.WriteJSON(gin.H{
			"type":  "error",
			"error": err.Error(),
		})
		return
	}
} 