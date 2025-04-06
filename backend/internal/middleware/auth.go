package middleware

import (
	"net/http"
	"strings"

	"CHIPMUNK-T0T/ito_web_app/internal/auth"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "認証が必要です"})
			return
		}

		// Bearer tokenの形式をチェック
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "無効な認証形式です"})
			return
		}

		userID, err := auth.ValidateToken(parts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "無効なトークンです"})
			return
		}

		// 検証済みのユーザーIDをコンテキストに保存
		c.Set("user_id", userID)
		c.Next()
	}
}

func validateToken(token string) (uint, error) {
	// TODO: JWTの検証処理を実装
	return 0, nil
} 