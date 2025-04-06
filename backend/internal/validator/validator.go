package validator

import "github.com/go-playground/validator/v10"

var validate = validator.New()

type CreateRoomRequest struct {
	Name        string `json:"name" binding:"required,min=3,max=20"`
	MaxPlayers  int    `json:"max_players" binding:"required,min=2,max=10"`
	Description string `json:"description" binding:"max=200"`
	IsPrivate   bool   `json:"is_private"`
	Password    string `json:"password" binding:"required_if=IsPrivate true"`
}

type JoinRoomRequest struct {
	Password string `json:"password"`
}

func ValidateStruct(s interface{}) error {
	return validate.Struct(s)
} 