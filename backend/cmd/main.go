package main

import (
	"log"
	"os"
	"time"

	"CHIPMUNK-T0T/ito_web_app/internal/database"
	"CHIPMUNK-T0T/ito_web_app/internal/logger"
	"CHIPMUNK-T0T/ito_web_app/internal/middleware"
	"CHIPMUNK-T0T/ito_web_app/internal/repository"
	"CHIPMUNK-T0T/ito_web_app/internal/router"
	"CHIPMUNK-T0T/ito_web_app/internal/usecase"
	"CHIPMUNK-T0T/ito_web_app/internal/websock"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// 環境変数のロード
	if err := godotenv.Load("../backend/.env"); err != nil {
		log.Printf("Warning: .env ファイルが見つかりません: %v", err)
	}
	
	// データベースの初期化
	db, err := database.NewDB_SQLite()
	if err != nil {
		log.Fatalf("データベース初期化エラー: %v", err)
	}

	if err := database.AutoMigrate(db); err != nil {
		log.Fatalf("マイグレーションエラー: %v", err)
	}

	if err := database.SeedThemes(db); err != nil {
		log.Printf("Warning: テーマのシードに失敗しました: %v", err)
	}

	// リポジトリの初期化
	userRepo := repository.NewUserRepository(db)
	roomRepo := repository.NewRoomRepository(db)
	themeRepo := repository.NewThemeRepository(db)
	// gameRepo := repository.NewGameRepository(db) // メモリ管理のため不要

	// GameHub の初期化
	gameHub := websock.NewGameHub()
	go gameHub.Run()

	// UseCase の初期化
	userUseCase := usecase.NewUserUseCase(userRepo)
	roomUseCase := usecase.NewRoomUseCase(roomRepo, userRepo)
	themeUseCase := usecase.NewThemeUseCase(themeRepo)
	gameUseCase := usecase.NewGameUseCase(roomUseCase, userUseCase, themeUseCase, gameHub)

	// WebSocket メッセージハンドラーの登録
	messageHandler := websock.NewMessageHandler(gameUseCase)
	gameHub.RegisterMessageHandler(messageHandler)

	// ロガーの初期化
	appLogger, err := logger.NewFileLogger("app.log")
	if err != nil {
		log.Fatalf("ロガー初期化エラー: %v", err)
	}
	defer appLogger.Close()

	// レートリミッターの初期化
	rateLimiter := middleware.NewRateLimiter()

	// Ginエンジンの初期化
	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.Use(rateLimiter.RateLimit(1000, time.Minute))

	// ルーターの初期化
	appRouter := router.NewRouter(userUseCase, roomUseCase, gameUseCase, appLogger)
	appRouter.SetupRoutes(engine)

	// サーバー起動
	if err := engine.Run(":8080"); err != nil {
		log.Fatalf("サーバー起動エラー: %v", err)
	}
}
