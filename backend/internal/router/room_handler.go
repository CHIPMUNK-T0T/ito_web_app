package router

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (r *Router) GetRooms(c *gin.Context) {
	rooms, err := r.roomUseCase.GetAllRooms()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ルーム一覧の取得に失敗しました"})
		return
	}

	c.JSON(http.StatusOK, rooms)
}

func (r *Router) CreateRoom(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		MaxPlayers  int    `json:"max_players" binding:"required,min=2,max=10"`
		Description string `json:"description"`
		IsPrivate   bool   `json:"is_private"`
		Password    string `json:"password"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無効なリクエスト形式です"})
		return
	}

	userID := c.GetUint("user_id")
	room, err := r.roomUseCase.CreateRoom(req.Name, req.Password, req.MaxPlayers, userID, req.Description, req.IsPrivate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, room)
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

	c.JSON(http.StatusOK, room)
}

func (r *Router) JoinRoom(c *gin.Context) {
	roomID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無効なルームIDです"})
		return
	}

	var req struct {
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無効なリクエスト形式です"})
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

	c.JSON(http.StatusOK, gin.H{
		"message": "ルームに参加しました",
		"room":    room,
	})
} 