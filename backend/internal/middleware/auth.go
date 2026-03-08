package middleware

import (
	"log"
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
			log.Printf("[Auth] Token validation failed: %v", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "無効なトークンです"})
			return
		}

		log.Printf("[Auth] Token validated. Path: %s, UserID: %d", c.Request.URL.Path, userID)

		// 検証済みのユーザーIDをコンテキストに保存
		c.Set("user_id", userID)
		c.Next()
	}
}

// AuthMiddlewareWS はWebSocket接続用の認証ミドルウェア。
// ブラウザのWebSocket APIはAuthorizationヘッダーを送れないため、
// クエリパラメータ ?token=<jwt> からトークンを読み取る。
// ヘッダーが存在する場合はAuthMiddlewareと同様にそちらを優先する。
func AuthMiddlewareWS() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := ""

		// 1. まずAuthorizationヘッダーを確認（通常HTTPリクエストとの互換性）
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			parts := strings.Split(authHeader, " ")
			if len(parts) == 2 && parts[0] == "Bearer" {
				tokenStr = parts[1]
			}
		}

		// 2. ヘッダーがなければクエリパラメータ ?token= を確認
		if tokenStr == "" {
			tokenStr = c.Query("token")
		}

		log.Printf("[AuthWS] Incoming request to %s with token: %s", c.Request.URL.Path, tokenStr)

		if tokenStr == "" {
			log.Printf("[AuthWS] No token found")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "認証が必要です"})
			return
		}

		userID, err := auth.ValidateToken(tokenStr)
		if err != nil {
			log.Printf("[AuthWS] Token validation failed: %v", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "無効なトークンです"})
			return
		}

		log.Printf("[AuthWS] Token validated. UserID: %d", userID)

		c.Set("user_id", userID)
		c.Next()
	}
}
