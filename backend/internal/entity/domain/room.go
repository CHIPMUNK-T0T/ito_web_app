package domain

import (
	"CHIPMUNK-T0T/ito_web_app/internal/entity"
	"CHIPMUNK-T0T/ito_web_app/internal/functional"
)

type Room struct {
	id         uint
	name       string
	password   functional.Hash
	maxPlayers int
	creatorID  uint

	description string
	isPrivate   bool
	status      entity.RoomStatus
	players     []RoomPlayer
}

func NewRoom(name, password string, maxPlayers int, creatorID uint, description string, isPrivate bool) (Room, error) {
	hash, err := functional.Encrypt(password)
	if err != nil {
		return Room{}, err

	}

	return Room{
		name:        name,
		password:    hash,
		maxPlayers:  maxPlayers,
		creatorID:   creatorID,
		description: description,
		isPrivate:   isPrivate,

		status:  entity.RoomStatusWaiting,
		players: make([]RoomPlayer, 0),
	}, nil
}

func NewRoomWithID(id uint, name string, passwordHash functional.Hash, maxPlayers int, creatorID uint, description string, isPrivate bool) Room {
	return Room{
		id:          id,
		name:        name,
		password:    passwordHash,
		maxPlayers:  maxPlayers,
		creatorID:   creatorID,
		description: description,

		isPrivate: isPrivate,

		status:  entity.RoomStatusWaiting,
		players: make([]RoomPlayer, 0),
	}
}

func (r Room) ID() uint {
	return r.id
}

func (r *Room) SetID(id uint) {
	r.id = id
}

func (r Room) Name() string {
	return r.name
}

func (r Room) Password() functional.Hash {
	return r.password
}

func (r Room) MaxPlayers() int {
	return r.maxPlayers
}

func (r Room) CreatorID() uint {
	return r.creatorID
}

func (r Room) Description() string {
	return r.description
}

func (r Room) IsPrivate() bool {
	return r.isPrivate
}

func (r Room) IsFull() bool {
	return len(r.players) >= r.maxPlayers
}

func (r Room) GetPlayerCount() int {
	return len(r.players)
}

func (r *Room) LoadPlayers(players []RoomPlayer) {
	r.players = players
}

func (r *Room) GetPlayers() []*RoomPlayer {
	players := make([]*RoomPlayer, len(r.players))
	for i := range r.players {
		players[i] = &r.players[i]
	}
	return players
}

func (r Room) GetStatus() entity.RoomStatus {
	return r.status
}

func (r Room) ValidatePassword(password string) bool {
	return r.password.Validate(password)
}

func (r *Room) AddPlayer(user User) bool {
	if len(r.players) >= r.maxPlayers {
		return false
	}

	for _, player := range r.players {
		if player.user.ID() == user.ID() {
			return false
		}
	}

	role := entity.Guest
	if r.creatorID == user.ID() {
		role = entity.Host

	}
	r.players = append(r.players, RoomPlayer{
		user:    user,
		role:    role,
		isReady: false,
	})

	return true

}

func (r *Room) RemovePlayer(userID uint) bool {
	index := -1
	for i, player := range r.players {
		if player.user.ID() == userID {
			index = i
			break
		}
	}

	if index == -1 {
		return false
	}

	// 削除実行
	r.players = append(r.players[:index], r.players[index+1:]...)

	// もし削除されたのがホスト（creatorID）だった場合、新しいホストを任命
	if r.creatorID == userID && len(r.players) > 0 {
		newHost := &r.players[0]
		r.creatorID = newHost.user.ID()
		newHost.role = entity.Host
	}

	return true
}

func (r *Room) GetCreatorID() uint {
	return r.creatorID
}

func (r *Room) HasPlayer(userID uint) bool {
	for _, player := range r.players {
		user := player.User()
		if user.ID() == userID {
			return true
		}
	}
	return false
}

type RoomPlayer struct {
	user        User
	role        entity.Role
	isReady     bool
	isConnected bool
	cardValue   *int // ゲーム開始時に配布される1-100の数字
}

func NewRoomPlayer(user User, role entity.Role) RoomPlayer {
	return RoomPlayer{
		user:      user,
		role:      role,
		isReady:   false,
	}
}

func (r RoomPlayer) User() User {
	return r.user
}

func (r *RoomPlayer) Role() entity.Role {
	return r.role
}

func (r *RoomPlayer) IsReady() bool {
	return r.isReady
}

func (r *RoomPlayer) CardValue() *int {
	return r.cardValue
}

func (r *RoomPlayer) SetRole(role entity.Role) {
	r.role = role
}

func (r *RoomPlayer) SetCardValue(cardValue int) {
	r.cardValue = &cardValue
}

func (r *RoomPlayer) SetIsReady(isReady bool) {
	r.isReady = isReady
}

func (r *RoomPlayer) IsHost() bool {
	return r.role == entity.Host
}

func (r *RoomPlayer) IsGuest() bool {
	return r.role == entity.Guest
}

func (r *RoomPlayer) SetIsConnected(connected bool) {
	r.isConnected = connected
}

func (r *RoomPlayer) IsConnected() bool {
	return r.isConnected
}