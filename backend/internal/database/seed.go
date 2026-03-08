package database

import (
	"CHIPMUNK-T0T/ito_web_app/internal/entity/model"
	"gorm.io/gorm"
)

func SeedThemes(db *gorm.DB) error {
	themes := []string{
		"朝食の値段",
		"理想の気温",
		"動物の大きさ",
		"食べ物の辛さ",
		"欲しいものの重さ",
		"ゾンビに襲われた時の武器の頼もしさ",
		"無人島に持っていきたいもの",
		"初デートで行きたい場所の混雑度",
	}

	for _, content := range themes {
		var theme model.Theme
		if err := db.Where("content = ?", content).First(&theme).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				db.Create(&model.Theme{Content: content})
			}
		}
	}
	return nil
}
