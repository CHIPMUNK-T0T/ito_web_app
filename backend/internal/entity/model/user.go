package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	// Username と Password の組み合わせでユニークとする
	Username string `gorm:"type:text;size:16;not null;uniqueIndex:idx_username_password"`
	Password string `gorm:"type:text;not null;uniqueIndex:idx_username_password"`
}
