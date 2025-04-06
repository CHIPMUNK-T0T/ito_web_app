package repository

import (
	"CHIPMUNK-T0T/ito_web_app/internal/entity/domain"
	"CHIPMUNK-T0T/ito_web_app/internal/entity/model"
	"CHIPMUNK-T0T/ito_web_app/internal/functional"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) IUserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user domain.User) error {
	modelUser := &model.User{
		Username: user.Username(),
		Password: string(user.Password()),
	}
	result := r.db.Create(modelUser)
	if result.Error != nil {
		return fmt.Errorf("ユーザー作成エラー: %v", result.Error)
	}
	return nil
}

func (r *userRepository) FindByID(id uint) (domain.User, error) {
	var modelUser model.User
	result := r.db.First(&modelUser, id)
	if result.Error != nil {
		return domain.User{}, fmt.Errorf("ユーザー取得エラー: %v", result.Error)
	}
	return domain.NewUserWithID(modelUser.ID, modelUser.Username, functional.Hash(modelUser.Password)), nil
}

func (r *userRepository) FindByUserNameAndPassword(username string, password string) (domain.User, error) {
	var modelUser model.User
	result := r.db.Where("username = ? AND password = ?", username, password).First(&modelUser)
	if result.Error != nil {
		return domain.User{}, fmt.Errorf("ユーザー取得エラー: %v", result.Error)
	}
	return domain.NewUserWithID(modelUser.ID, modelUser.Username, functional.Hash(modelUser.Password)), nil
}

func (r *userRepository) FindAll() ([]domain.User, error) {
	var modelUsers []model.User
	result := r.db.Find(&modelUsers)
	if result.Error != nil {
		return nil, fmt.Errorf("ユーザー一覧取得エラー: %v", result.Error)
	}

	domainUsers := make([]domain.User, 0, len(modelUsers))
	for _, modelUser := range modelUsers {
		domainUser := domain.NewUserWithID(modelUser.ID, modelUser.Username, functional.Hash(modelUser.Password))
		domainUsers = append(domainUsers, domainUser)
	}
	return domainUsers, nil
}

func (r *userRepository) Update(user domain.User) error {
	var pastModel model.User
	result := r.db.First(&pastModel, user.ID())
	if result.Error != nil {
		return fmt.Errorf("ユーザー更新エラー: %v", result.Error)
	}

	modelUser := &model.User{
		Model:      gorm.Model{ID: pastModel.ID, UpdatedAt: time.Now(), CreatedAt: pastModel.CreatedAt, DeletedAt: pastModel.DeletedAt},
		Username:   user.Username(),
		Password:   string(user.Password()),
	}
	result = r.db.Save(modelUser)
	if result.Error != nil {
		return fmt.Errorf("ユーザー更新エラー: %v", result.Error)
	}
	return nil
}

func (r *userRepository) Delete(id uint) error {
	result := r.db.Delete(&model.User{}, id)
	if result.Error != nil {
		return fmt.Errorf("ユーザー削除エラー: %v", result.Error)
	}
	return nil
}
