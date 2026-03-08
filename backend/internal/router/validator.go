package router

import (
	"CHIPMUNK-T0T/ito_web_app/internal/validator"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	goValidator "github.com/go-playground/validator/v10"
)

func getErrorMsg(fe goValidator.FieldError) string {
	field := fe.Field()
	// 簡単な日本語マッピング
	switch field {
	case "Name":
		field = "ルーム名"
	case "MaxPlayers":
		field = "最大参加人数"
	case "Description":
		field = "説明"
	case "Password":
		field = "パスワード"
	}

	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("%sは必須です。", field)
	case "min":
		return fmt.Sprintf("%sは%s文字以上で入力してください。", field, fe.Param())
	case "max":
		return fmt.Sprintf("%sは%s文字以内で入力してください。", field, fe.Param())
	case "required_if":
		// required_if=IsPrivate true
		parts := strings.Split(fe.Param(), " ")
		if len(parts) > 0 && parts[0] == "IsPrivate" {
			return fmt.Sprintf("プライベートルームにする場合、%sは必須です。", field)
		}
		return fmt.Sprintf("%sは特定の条件下で必須です。", field)
	default:
		return fmt.Sprintf("%sでエラーが発生しました。", fe.Field())
	}
}

func (r *Router) validateCreateRoom(c *gin.Context) {
	var req validator.CreateRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		var ve goValidator.ValidationErrors
		if errors.As(err, &ve) {
			out := make([]string, len(ve))
			for i, fe := range ve {
				out[i] = getErrorMsg(fe)
			}
			c.JSON(http.StatusBadRequest, gin.H{"error": strings.Join(out, "\n")})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "不正なリクエストです。"})
		}
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
 