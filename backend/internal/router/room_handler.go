package router

import (
	"CHIPMUNK-T0T/ito_web_app/internal/entity/domain"
	"CHIPMUNK-T0T/ito_web_app/internal/validator"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type playerResponse struct {
	ID        uint   `json:"id"`
	Username  string `json:"username"`
	Role      string `json:"role"`
	IsReady   bool   `json:"isReady"`
	IsConnect bool   `json:"isConnected"`
}

type roomResponse struct {
	ID          uint             `json:"id"`
	Name        string           `json:"name"`
	MaxPlayers  int              `json:"max_players"`
	Description string           `json:"description,omitempty"`
	IsPrivate   bool             `json:"is_private"`
	CreatorID   uint             `json:"creator_id"`
	PlayerCount int              `json:"player_count"`
	Players     []playerResponse `json:"players"`
}

func newRoomResponse(room *domain.Room) roomResponse {
	if room == nil {
		return roomResponse{}
	}

	players := room.GetPlayers()
	playerResponses := make([]playerResponse, 0, len(players))
	for _, p := range players {
		playerResponses = append(playerResponses, playerResponse{
			ID:        p.User().ID(),
			Username:  p.User().Username(),
			Role:      string(p.Role()),
			IsReady:   p.IsReady(),
			IsConnect: p.IsConnected(),
		})
	}

	return roomResponse{
		ID:          room.ID(),
		Name:        room.Name(),
		MaxPlayers:  room.MaxPlayers(),
		Description: room.Description(),
		IsPrivate:   room.IsPrivate(),
		CreatorID:   room.CreatorID(),
		PlayerCount: room.GetPlayerCount(),
		Players:     playerResponses,
	}
}

func roomsToResponse(rooms []domain.Room) []roomResponse {
	responses := make([]roomResponse, 0, len(rooms))
	for i := range rooms {
		responses = append(responses, newRoomResponse(&rooms[i]))
	}
	return responses
}

func (r *Router) GetRooms(c *gin.Context) {
	rooms, err := r.roomUseCase.GetAllRooms()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ルーム一覧の取得に失敗しました"})
		return
	}

	c.JSON(http.StatusOK, roomsToResponse(rooms))
}

func (r *Router) CreateRoom(c *gin.Context) {
	reqData, exists := c.Get("createRoomRequest")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "リクエストデータの取得に失敗しました"})
		return
	}

	req, ok := reqData.(validator.CreateRoomRequest)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "リクエストデータの型が不正です"})
		return
	}

	userID := c.GetUint("user_id")
	room, err := r.roomUseCase.CreateRoom(req.Name, req.Password, req.MaxPlayers, userID, req.Description, req.IsPrivate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// ルーム作成者（ホスト）をプレイヤーとして追加する
	if err := r.roomUseCase.AddPlayer(room.ID(), userID); err != nil {
		r.logger.Error("ホストのAddPlayerに失敗しました", "roomID", room.ID(), "userID", userID, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ルーム作成後のプレイヤー追加に失敗しました"})
		return
	}

	// AddPlayer後にルームを再取得して最新のプレイヤー情報を返す
	updatedRoom, err := r.roomUseCase.GetRoom(room.ID())
	if err != nil {
		c.JSON(http.StatusCreated, newRoomResponse(room))
		return
	}
	c.JSON(http.StatusCreated, newRoomResponse(updatedRoom))
}

func (r *Router) GetRoom(c *gin.Context) {
	roomID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無効なルームIDです"})
		return
	}

	room, err := r.roomUseCase.GetRoom(uint(roomID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ルームが見つかりません"})
		return
	}

	c.JSON(http.StatusOK, newRoomResponse(room))
}

func (r *Router) JoinRoom(c *gin.Context) {
	roomIDStr := c.Param("id")
	roomID, err := strconv.ParseUint(roomIDStr, 10, 64)
	if err != nil {
		fmt.Printf("[Handler] Invalid RoomID param: %s\n", roomIDStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "無効なルームIDです"})
		return
	}

	userID := c.GetUint("user_id")
	fmt.Printf("[Handler] User %d START JoinRoom for ID %d\n", userID, roomID)

	reqData, exists := c.Get("joinRoomRequest")
	if !exists {
		fmt.Printf("[Handler] joinRoomRequest missing in context for User %d\n", userID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "リクエストデータの取得に失敗しました"})
		return
	}
	req, ok := reqData.(validator.JoinRoomRequest)
	if !ok {
		fmt.Printf("[Handler] joinRoomRequest type mismatch for User %d\n", userID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "リクエストデータの型が不正です"})
		return
	}

	room, err := r.roomUseCase.JoinRoom(uint(roomID), req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetUint("user_id")
	if err := r.roomUseCase.AddPlayer(uint(roomID), userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 最新の状態を取得して返す
	updatedRoom, err := r.roomUseCase.GetRoom(uint(roomID))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "ルームに参加しました",
			"room":    newRoomResponse(room),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "ルームに参加しました",
		"room":    newRoomResponse(updatedRoom),
	})
}
