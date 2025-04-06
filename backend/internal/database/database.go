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

	// // 既存のテーブルを削除（開発環境でのみ使用することを推奨）
	// err = db.Migrator().DropTable(&model.User{}, &model.Room{}, &model.RoomUser{})
	// if err != nil {
	// 	return nil, fmt.Errorf("テーブル削除エラー: %v", err)
	// }

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("コネクションプール取得エラー: %v", err)
	}

	// コネクションプールの設定
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	log.Println("データベース接続成功")
	return db, nil
}

func AutoMigrate(db *gorm.DB) error {
	// モデルの自動マイグレーション
	return db.AutoMigrate(
		&model.User{},
		&model.Room{},
	)
}

func GetTableNames(db *gorm.DB) ([]string, error) {
	tables, err := db.Migrator().GetTables()
	if err != nil {
		return nil, fmt.Errorf("テーブル名取得エラー: %v", err)
	}
	return tables, nil
}

