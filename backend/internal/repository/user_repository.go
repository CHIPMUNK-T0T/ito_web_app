package repository

import "CHIPMUNK-T0T/ito_web_app/internal/entity/domain"

type IUserRepository interface {
	Create(user domain.User) error
	FindByID(id uint) (domain.User, error)
	FindByUserNameAndPassword(username string, password string) (domain.User, error)
	FindAll() ([]domain.User, error)
	Update(user domain.User) error
	Delete(id uint) error
}
