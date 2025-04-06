package model

import (
	"gorm.io/gorm"
)

type Room struct {
	gorm.Model
	Name        string     `gorm:"type:text;size:16;not null;uniqueIndex:idx_name_password"`
	Password    string     `gorm:"type:text;size:16;not null;uniqueIndex:idx_name_password"`
	MaxPlayers  int        `gorm:"not null"`
	CreatorID   uint       `gorm:"not null"`
	Description string     `gorm:"type:text;size:300"`
	IsPrivate   bool       `gorm:"default:false"`
}
