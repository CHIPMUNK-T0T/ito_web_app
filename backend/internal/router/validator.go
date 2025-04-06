package router

import (
	"CHIPMUNK-T0T/ito_web_app/internal/validator"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (r *Router) validateCreateRoom(c *gin.Context) {
	var req validator.CreateRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "入力値が不正です"})
		c.Abort()
		return
	}

	c.Set("createRoomRequest", req)
	c.Next()
}

func (r *Router) validateJoinRoom(c *gin.Context) {
	var req validator.JoinRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "入力値が不正です"})
		c.Abort()
		return
	}

	c.Set("joinRoomRequest", req)
	c.Next()
} 