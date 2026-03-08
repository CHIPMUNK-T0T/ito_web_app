package usecase

import (
	"errors"

	"CHIPMUNK-T0T/ito_web_app/internal/entity/domain"
	"CHIPMUNK-T0T/ito_web_app/internal/repository"
)

type UserUseCase struct {
	userRepo repository.IUserRepository
}

func NewUserUseCase(userRepo repository.IUserRepository) *UserUseCase {
	return &UserUseCase{
		userRepo: userRepo,
	}
}

func (uc *UserUseCase) Register(username, password string) (*domain.User, error) {
	if len(username) < 3 || len(username) > 16 {
		return nil, errors.New("ユーザー名は3文字以上16文字以下である必要があります")
	}

	if len(password) < 8 || len(password) > 16 {
		return nil, errors.New("パスワードは8文字以上16文字以下である必要があります")
	}

	user, err := domain.NewUser(username, password)
	if err != nil {
		return nil, err
	}

	err = uc.userRepo.Create(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (uc *UserUseCase) Login(username, password string) (*domain.User, error) {
	user, err := uc.userRepo.FindByUserNameAndPassword(username, password)
	if err != nil {
		return nil, errors.New("ユーザー名またはパスワードが正しくありません")
	}

	return &user, nil
}

func (uc *UserUseCase) GetUserByID(id uint) (*domain.User, error) {
	user, err := uc.userRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
