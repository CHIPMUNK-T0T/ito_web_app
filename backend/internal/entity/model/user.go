package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"type:text;size:16;not null;uniqueIndex:idx_username_password"`
	Password string `gorm:"type:text;size:16;not null;uniqueIndex:idx_username_password"`
}
