package usecase

import (
	"errors"

	"CHIPMUNK-T0T/ito_web_app/internal/entity/domain"
	"CHIPMUNK-T0T/ito_web_app/internal/repository"
)

type RoomUseCase struct {
	roomRepo repository.IRoomRepository
	userRepo repository.IUserRepository
}

func NewRoomUseCase(roomRepo repository.IRoomRepository, userRepo repository.IUserRepository) *RoomUseCase {
	return &RoomUseCase{
		roomRepo: roomRepo,
		userRepo: userRepo,
	}
}

func (uc *RoomUseCase) CreateRoom(name, password string, maxPlayers int, creatorID uint, description string, isPrivate bool) (*domain.Room, error) {
	if len(name) < 3 || len(name) > 16 {
		return nil, errors.New("ルーム名は3文字以上16文字以下である必要があります")
	}

	if maxPlayers < 2 || maxPlayers > 10 {
		return nil, errors.New("プレイヤー数は2人以上10人以下である必要があります")
	}

	if len(description) > 300 {
		return nil, errors.New("説明文は300文字以内である必要があります")
	}

	room, err := domain.NewRoom(name, password, maxPlayers, creatorID, description, isPrivate)
	if err != nil {
		return nil, err
	}

	err = uc.roomRepo.Create(room)
	if err != nil {
		return nil, err
	}

	return &room, nil
}

func (uc *RoomUseCase) JoinRoom(roomID uint, password string) (*domain.Room, error) {
	room, err := uc.roomRepo.FindByID(roomID)
	if err != nil {
		return nil, errors.New("ルームが見つかりません")
	}

	if !room.ValidatePassword(password) {
		return nil, errors.New("パスワードが一致しません")
	}

	return &room, nil
}

func (uc *RoomUseCase) AddPlayer(roomID, userID uint) error {
	room, err := uc.roomRepo.FindByID(roomID)
	if err != nil {
		return errors.New("ルームが見つかりません")
	}

	user, err := uc.userRepo.FindByID(userID)
	if err != nil {
		return errors.New("ユーザーが見つかりません")
	}

	if room.IsFull() {
		return errors.New("ルームが満員です")
	}

	if !room.AddPlayer(user) {
		return errors.New("プレイヤーの追加に失敗しました")
	}

	return uc.roomRepo.Update(room)
}

func (uc *RoomUseCase) RemovePlayer(roomID uint, userID uint) error {
	room, err := uc.roomRepo.FindByID(roomID)
	if err != nil {
		return errors.New("ルームが見つかりません")
	}

	if !room.RemovePlayer(userID) {
		return errors.New("プレイヤーの削除に失敗しました")
	}

	return uc.roomRepo.Update(room)
}

func (uc *RoomUseCase) SetPlayerReady(roomID uint, userID uint, isReady bool) error {
	room, err := uc.roomRepo.FindByID(roomID)
	if err != nil {
		return errors.New("ルームが見つかりません")
	}

	players := room.GetPlayers()
	for i := range players {
		user := players[i].User()
		if user.ID() == userID {
			players[i].SetIsReady(isReady)
			return uc.roomRepo.Update(room)
		}
	}

	return errors.New("プレイヤーが見つかりません")
}

func (uc *RoomUseCase) GetRoom(roomID uint) (*domain.Room, error) {
	room, err := uc.roomRepo.FindByID(roomID)
	if err != nil {
		return nil, errors.New("ルームが見つかりません")
	}
	return &room, nil
}

func (uc *RoomUseCase) GetAllRooms() ([]domain.Room, error) {
	return uc.roomRepo.FindAll()
} 