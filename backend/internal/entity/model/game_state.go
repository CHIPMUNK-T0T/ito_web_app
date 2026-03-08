package model

import (
	"gorm.io/gorm"
)

// GameStatus represents the current stage of a game session
type GameStatus string

const (
	StatusWaiting  GameStatus = "waiting"
	StatusPlaying  GameStatus = "playing"
	StatusFinished GameStatus = "finished"
)

// GameSession tracks an active or past game in a room
type GameSession struct {
	gorm.Model
	RoomID       uint       `gorm:"not null;index"`
	Status       GameStatus `gorm:"type:text;not null;default:'waiting'"`
	ThemeContent string     `gorm:"type:text"`
	CurrentRound int        `gorm:"not null;default:1"`
}

// Hand tracks the cards held by a player in a specific game session
type Hand struct {
	gorm.Model
	GameSessionID uint `gorm:"not null;index"`
	UserID        uint `gorm:"not null;index"`
	CardValue     int  `gorm:"not null"`
	IsRevealed    bool `gorm:"not null;default:false"`
}

// Theme represents the available categories for the game
type Theme struct {
	gorm.Model
	Content string `gorm:"type:text;not null;unique"`
}
