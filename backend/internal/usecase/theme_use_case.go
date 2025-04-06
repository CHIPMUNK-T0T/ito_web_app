package usecase

import "CHIPMUNK-T0T/ito_web_app/internal/repository"

type ThemeUseCase struct {
	themeRepo repository.IThemeRepository
}

func NewThemeUseCase(themeRepo repository.IThemeRepository) *ThemeUseCase {
	return &ThemeUseCase{
		themeRepo: themeRepo,
	}
}

func (uc *ThemeUseCase) GetRandomTheme() (string, error) {
	return uc.themeRepo.GetRandom()
}
