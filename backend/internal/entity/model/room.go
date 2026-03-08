package model

import (
	"gorm.io/gorm"
)

type Room struct {
	gorm.Model
	Name        string `gorm:"type:text;not null"`
	Password    string `gorm:"type:text;not null"` // bcryptハッシュを格納するためsizeを制限しない
	MaxPlayers  int    `gorm:"not null"`
	CreatorID   uint   `gorm:"not null"`
	Description string `gorm:"type:text"`
	IsPrivate   bool   `gorm:"default:false"`
}
