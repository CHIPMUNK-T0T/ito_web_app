package repository

type IThemeRepository interface {
	FindAll() ([]string, error)
	Create(theme string) error
	GetRandom() (string, error)
}