package repository

import "CHIPMUNK-T0T/ito_web_app/internal/entity/domain"

type IRoomRepository interface {
	Create(room domain.Room) error
	FindByID(id uint) (domain.Room, error)
	FindByRoomNameAndPassword(name string, password string) (domain.Room, error)
	FindAll() ([]domain.Room, error)
	Update(room domain.Room) error
	Delete(id uint) error
}
