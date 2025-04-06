package repository

import (
	"gorm.io/gorm"
)

type ThemeRepository struct {
	db *gorm.DB
}

func NewThemeRepository(db *gorm.DB) *ThemeRepository {
	return &ThemeRepository{db: db}
}

func (r *ThemeRepository) FindAll() ([]string, error) {
	var themes []string
	if err := r.db.Table("themes").Pluck("content", &themes).Error; err != nil {
		return nil, err
	}
	return themes, nil
}

func (r *ThemeRepository) Create(theme string) error {
	return r.db.Table("themes").Create(map[string]interface{}{
		"content": theme,
	}).Error
}

func (r *ThemeRepository) GetRandom() (string, error) {
	var theme string
	if err := r.db.Table("themes").Order("RANDOM()").Limit(1).Pluck("content", &theme).Error; err != nil {
		return "", err
	}
	return theme, nil
}
