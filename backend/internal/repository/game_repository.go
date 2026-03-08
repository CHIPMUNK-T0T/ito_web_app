package repository

import (
	"CHIPMUNK-T0T/ito_web_app/internal/entity/model"
	"gorm.io/gorm"
)

type GameRepository struct {
	db *gorm.DB
}

func NewGameRepository(db *gorm.DB) *GameRepository {
	return &GameRepository{db: db}
}

func (r *GameRepository) CreateSession(session *model.GameSession) error {
	return r.db.Create(session).Error
}

func (r *GameRepository) UpdateSession(session *model.GameSession) error {
	return r.db.Save(session).Error
}

func (r *GameRepository) GetSessionByRoomID(roomID uint) (*model.GameSession, error) {
	var session model.GameSession
	err := r.db.Where("room_id = ? AND status != ?", roomID, model.StatusFinished).Last(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *GameRepository) SaveHand(hand *model.Hand) error {
	return r.db.Save(hand).Error
}

func (r *GameRepository) GetHandsBySessionID(sessionID uint) ([]model.Hand, error) {
	var hands []model.Hand
	err := r.db.Where("game_session_id = ?", sessionID).Find(&hands).Error
	return hands, err
}

func (r *GameRepository) GetHandByUser(sessionID, userID uint) (*model.Hand, error) {
	var hand model.Hand
	err := r.db.Where("game_session_id = ? AND user_id = ?", sessionID, userID).First(&hand).Error
	if err != nil {
		return nil, err
	}
	return &hand, nil
}

func (r *GameRepository) Transaction(fn func(repo *GameRepository) error) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		newRepo := &GameRepository{db: tx}
		return fn(newRepo)
	})
}

