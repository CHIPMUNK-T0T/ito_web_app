package main

import (
	"CHIPMUNK-T0T/ito_web_app/internal/database"
	"CHIPMUNK-T0T/ito_web_app/internal/repository"
	"log"

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

	// マイグレーションの実行
	// if err := database.AutoMigrate(db); err != nil {
	// 	log.Fatalf("マイグレーションエラー: %v", err)
	// }

	return

	// リポジトリの初期化
	// userRepo := repository.NewUserRepository(db)
	// roomRepo := repository.NewRoomRepository(db)
	// themeRepo := repository.NewThemeRepository(db)

	// // GameHubの初期化
	// gameHub := websock.NewGameHub()
	// go gameHub.Run()

	// // UseCaseの初期化
	// userUseCase := usecase.NewUserUseCase(userRepo)
	// roomUseCase := usecase.NewRoomUseCase(roomRepo, userRepo)
	// themeUseCase := usecase.NewThemeUseCase(themeRepo)
	// gameUseCase := usecase.NewGameUseCase(roomUseCase, userUseCase, themeUseCase, gameHub)

	// // WebSocketハンドラーの作成
	// // messageHandler := websock.NewMessageHandler(gameUseCase)
	// // gameHub.RegisterMessageHandler(messageHandler)

	// // 宣言エラー回避
	// fmt.Println(no_use(userRepo, roomRepo))

	// // ロガーの初期化
	// appLogger, err := logger.NewFileLogger("app.log")
	// if err != nil {
	// 	log.Fatalf("ロガー初期化エラー: %v", err)
	// }
	// defer appLogger.Close()

	// // レートリミッターの初期化
	// rateLimiter := middleware.NewRateLimiter()

	// // Ginエンジンの初期化
	// engine := gin.Default()

	// // ミドルウェアの設定
	// engine.Use(middleware.LoggerMiddleware(appLogger))
	// engine.Use(middleware.CorsMiddleware())
	// engine.Use(rateLimiter.RateLimit(100, time.Minute)) // 1分間に100リクエストまで

	// // Routerの初期化
	// router := router.NewRouter(userUseCase, roomUseCase, gameUseCase, appLogger)
	
	// // ルートの設定
	// router.SetupRoutes(engine)

	// // サーバー起動
	// if err := engine.Run(":8080"); err != nil {
	// 	log.Fatalf("サーバー起動エラー: %v", err)
	// }
}

func no_use(IUserRepository repository.IUserRepository, IRoomRepository repository.IRoomRepository) string {
	return "test"
}
