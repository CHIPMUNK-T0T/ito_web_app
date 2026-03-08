package database

import (
	"CHIPMUNK-T0T/ito_web_app/internal/entity/model"
	"fmt"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewDB_SQLite() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("ito_game.db"), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("データベース接続エラー: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("コネクションプール取得エラー: %v", err)
	}

	// コネクションプールの設定
	// SQLite はデフォルトで単一書き込みのため、1 に制限することでロック競合を防ぐ
	sqlDB.SetMaxIdleConns(1)
	sqlDB.SetMaxOpenConns(1)

	log.Println("データベース接続成功")
	return db, nil
}

func AutoMigrate(db *gorm.DB) error {
	// モデルの自動マイグレーション
	// ゲームの進行状況（GameSession, Hand）はメモリ管理にするため、DBからは除外
	return db.AutoMigrate(
		&model.User{},
		&model.Room{},
		&model.Theme{},
	)
}

func GetTableNames(db *gorm.DB) ([]string, error) {
	tables, err := db.Migrator().GetTables()
	if err != nil {
		return nil, fmt.Errorf("テーブル名取得エラー: %v", err)
	}
	return tables, nil
}
