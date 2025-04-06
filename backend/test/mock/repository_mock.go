package mock

import (
	"CHIPMUNK-T0T/ito_web_app/internal/entity/domain"
	"CHIPMUNK-T0T/ito_web_app/internal/repository"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewMockDB() (*gorm.DB, sqlmock.Sqlmock, error) {
	// sqlmock でモックの *sql.DB を生成
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, fmt.Errorf("unexpected error when opening a stub database connection: %s", err)
	}

	// MySQL Dialector を利用して *sql.DB をラップ
	dialector := mysql.New(mysql.Config{
		Conn:                      sqlDB,
		SkipInitializeWithVersion: true, // バージョンチェックをスキップ
	})

	// GORM の DB インスタンスを生成
	gormDB, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open gorm DB: %v", err)
	}

	return gormDB, mock, nil
}

type UserRepository struct {
	db   *gorm.DB
	mock sqlmock.Sqlmock
}

func NewUserRepository(db *gorm.DB, mock sqlmock.Sqlmock) repository.IUserRepository {
	return &UserRepository{
		db:   db,
		mock: mock,
	}
}

func (r *UserRepository) Create(user domain.User) error {
	result := r.db.Create(&user)
	if result.Error != nil {
		return fmt.Errorf("failed to create user: %v", result.Error)
	}
	return nil
}

func (r *UserRepository) FindByID(id uint) (domain.User, error) {
	var user domain.User
	result := r.db.First(&user, id)
	if result.Error != nil {
		return domain.User{}, fmt.Errorf("user not found with id: %d", id)
	}
	return user, nil
}

func (r *UserRepository) FindByUserNameAndPassword(username string, password string) (domain.User, error) {
	var user domain.User
	result := r.db.Where("username = ?", username).First(&user)
	if result.Error != nil {
		return domain.User{}, fmt.Errorf("user not found: %s", username)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password()), []byte(password)); err != nil {
		return domain.User{}, fmt.Errorf("invalid password for user: %s", username)
	}
	return user, nil
}

func (r *UserRepository) FindAll() ([]domain.User, error) {
	var users []domain.User
	result := r.db.Find(&users)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to fetch users: %v", result.Error)
	}
	return users, nil
}

func (r *UserRepository) Update(user domain.User) error {
	result := r.db.Save(&user)
	if result.Error != nil {
		return fmt.Errorf("failed to update user: %v", result.Error)
	}
	return nil
}

func (r *UserRepository) Delete(id uint) error {
	result := r.db.Delete(&domain.User{}, id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete user: %v", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("user not found with id: %d", id)
	}
	return nil
}

func (r *UserRepository) Clear() {
	// SQLモックを使用する場合はClearは不要
}

type RoomRepository struct {
	db   *gorm.DB
	mock sqlmock.Sqlmock
}

func NewRoomRepository(db *gorm.DB, mock sqlmock.Sqlmock) repository.IRoomRepository {
	return &RoomRepository{
		db:   db,
		mock: mock,
	}
}

func (r *RoomRepository) Create(room domain.Room) error {
	result := r.db.Create(&room)
	if result.Error != nil {
		return fmt.Errorf("failed to create room: %v", result.Error)
	}
	return nil
}

func (r *RoomRepository) FindByID(id uint) (domain.Room, error) {
	var room domain.Room
	result := r.db.First(&room, id)
	if result.Error != nil {
		return domain.Room{}, fmt.Errorf("room not found with id: %d", id)
	}
	return room, nil
}

func (r *RoomRepository) FindByRoomNameAndPassword(name, password string) (domain.Room, error) {
	var room domain.Room
	result := r.db.Where("name = ?", name).First(&room)
	if result.Error != nil {
		return domain.Room{}, fmt.Errorf("room not found with name: %s", name)
	}
	if !room.ValidatePassword(password) {
		return domain.Room{}, fmt.Errorf("invalid password for room: %s", name)
	}
	return room, nil
}

func (r *RoomRepository) FindAll() ([]domain.Room, error) {
	var rooms []domain.Room
	result := r.db.Find(&rooms)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to fetch rooms: %v", result.Error)
	}
	return rooms, nil
}

func (r *RoomRepository) Update(room domain.Room) error {
	result := r.db.Save(&room)
	if result.Error != nil {
		return fmt.Errorf("failed to update room: %v", result.Error)
	}
	return nil
}

func (r *RoomRepository) Delete(id uint) error {
	result := r.db.Delete(&domain.Room{}, id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete room: %v", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("room not found with id: %d", id)
	}
	return nil
}

func (r *RoomRepository) Clear() {
	// SQLモックを使用する場合はClearは不要
}

type ThemeRepository struct {
	themes []string
}

func NewThemeRepository() repository.IThemeRepository {
	return &ThemeRepository{
		themes: []string{
			"一番安い買い物",
			"一番高い買い物", 
			"一番長く歩いた距離",
			"朝食の値段",
			"好きな数字",
			"自分の身長",
			"理想の気温",
			"理想の年収",
		},
	}
}

func (r *ThemeRepository) FindAll() ([]string, error) {
	if len(r.themes) == 0 {
		return nil, fmt.Errorf("no themes available")
	}
	return r.themes, nil
}

func (r *ThemeRepository) Create(theme string) error {
	r.themes = append(r.themes, theme)
	return nil
}

func (r *ThemeRepository) GetRandom() (string, error) {
	if len(r.themes) == 0 {
		return "", fmt.Errorf("no themes available")
	}
	return r.themes[0], nil
}

func (r *ThemeRepository) Clear() {
	r.themes = []string{}
}
