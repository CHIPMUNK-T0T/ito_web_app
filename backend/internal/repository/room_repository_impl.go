package repository

import (
	"CHIPMUNK-T0T/ito_web_app/internal/entity/domain"
	"CHIPMUNK-T0T/ito_web_app/internal/entity/model"
	"CHIPMUNK-T0T/ito_web_app/internal/functional"
	"fmt"
	"sync"
	"time"

	"gorm.io/gorm"
)

var (
	roomPlayersMap = make(map[uint][]domain.RoomPlayer)
	playersMu      sync.RWMutex
)

type roomRepository struct {
	db       *gorm.DB
	userRepo IUserRepository
}

func NewRoomRepository(db *gorm.DB) IRoomRepository {
	return &roomRepository{
		db:       db,
		userRepo: NewUserRepository(db),
	}
}

func (r *roomRepository) Create(room *domain.Room) error {
	modelRoom := &model.Room{
		Name:        room.Name(),
		Password:    string(room.Password()),
		MaxPlayers:  room.MaxPlayers(),
		CreatorID:   room.CreatorID(),
		Description: room.Description(),
		IsPrivate:   room.IsPrivate(),
	}
	result := r.db.Create(modelRoom)
	if result.Error != nil {
		return fmt.Errorf("ルーム作成エラー: %v", result.Error)
	}
	room.SetID(modelRoom.ID)

	playersMu.Lock()
	players := room.GetPlayers()
	saved := make([]domain.RoomPlayer, len(players))
	for i, p := range players {
		saved[i] = *p
	}
	roomPlayersMap[modelRoom.ID] = saved
	playersMu.Unlock()

	return nil
}

func (r *roomRepository) FindByID(id uint) (domain.Room, error) {
	var modelRoom model.Room
	result := r.db.First(&modelRoom, id)
	if result.Error != nil {
		return domain.Room{}, fmt.Errorf("ルーム取得エラー: %v", result.Error)
	}

	room := domain.NewRoomWithID(modelRoom.ID, modelRoom.Name, functional.Hash(modelRoom.Password), modelRoom.MaxPlayers, modelRoom.CreatorID, modelRoom.Description, modelRoom.IsPrivate)

	playersMu.RLock()
	if savedPlayers, ok := roomPlayersMap[modelRoom.ID]; ok {
		room.LoadPlayers(savedPlayers)
	}
	playersMu.RUnlock()

	return room, nil
}

func (r *roomRepository) FindByRoomNameAndPassword(name string, password string) (domain.Room, error) {
	var modelRoom model.Room
	result := r.db.Where("name = ? AND password = ?", name, password).First(&modelRoom)
	if result.Error != nil {
		return domain.Room{}, fmt.Errorf("ルーム取得エラー: %v", result.Error)
	}

	room := domain.NewRoomWithID(modelRoom.ID, modelRoom.Name, functional.Hash(modelRoom.Password), modelRoom.MaxPlayers, modelRoom.CreatorID, modelRoom.Description, modelRoom.IsPrivate)

	playersMu.RLock()
	if savedPlayers, ok := roomPlayersMap[modelRoom.ID]; ok {
		room.LoadPlayers(savedPlayers)
	}
	playersMu.RUnlock()

	return room, nil
}

func (r *roomRepository) FindAll() ([]domain.Room, error) {
	var modelRooms []model.Room
	result := r.db.Find(&modelRooms)
	if result.Error != nil {
		return nil, fmt.Errorf("ルーム一覧取得エラー: %v", result.Error)
	}

	domainRooms := make([]domain.Room, 0, len(modelRooms))
	playersMu.RLock()
	for _, modelRoom := range modelRooms {
		room := domain.NewRoomWithID(modelRoom.ID, modelRoom.Name, functional.Hash(modelRoom.Password), modelRoom.MaxPlayers, modelRoom.CreatorID, modelRoom.Description, modelRoom.IsPrivate)
		if savedPlayers, ok := roomPlayersMap[modelRoom.ID]; ok {
			room.LoadPlayers(savedPlayers)
		}
		domainRooms = append(domainRooms, room)
	}
	playersMu.RUnlock()
	return domainRooms, nil
}

func (r *roomRepository) Update(room *domain.Room) error {
	var pastModel model.Room
	result := r.db.First(&pastModel, room.ID())
	if result.Error != nil {
		return fmt.Errorf("ルーム更新エラー: %v", result.Error)
	}

	modelRoom := &model.Room{
		Model:       gorm.Model{ID: pastModel.ID, UpdatedAt: time.Now(), CreatedAt: pastModel.CreatedAt, DeletedAt: pastModel.DeletedAt},
		Name:        room.Name(),
		Password:    string(room.Password()),
		MaxPlayers:  room.MaxPlayers(),
		CreatorID:   room.CreatorID(),
		Description: room.Description(),
		IsPrivate:   room.IsPrivate(),
	}

	tx := r.db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("トランザクション開始エラー: %v", tx.Error)
	}

	if err := tx.Save(modelRoom).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("ルーム更新エラー: %v", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("トランザクションコミットエラー: %v", err)
	}

	playersMu.Lock()
	players := room.GetPlayers()
	saved := make([]domain.RoomPlayer, len(players))
	for i, p := range players {
		saved[i] = *p
	}
	roomPlayersMap[room.ID()] = saved
	playersMu.Unlock()

	return nil
}

func (r *roomRepository) Delete(id uint) error {
	tx := r.db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("トランザクション開始エラー: %v", tx.Error)
	}

	if err := tx.Delete(&model.Room{}, id).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("ルーム削除エラー: %v", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("トランザクションコミットエラー: %v", err)
	}

	playersMu.Lock()
	delete(roomPlayersMap, id)
	playersMu.Unlock()

	return nil
}
