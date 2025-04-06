package main

import (
	"CHIPMUNK-T0T/ito_web_app/internal/database"
	"CHIPMUNK-T0T/ito_web_app/internal/repository"
	"log"
)

func sub() {
	db, err := database.NewDB_SQLite()
	if err != nil {
		log.Fatalf("データベース初期化エラー: %v", err)
	}


	userRep := repository.NewUserRepository(db)


	return
}