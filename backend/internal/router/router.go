package router

import (
	"CHIPMUNK-T0T/ito_web_app/internal/logger"
	"CHIPMUNK-T0T/ito_web_app/internal/middleware"
	"CHIPMUNK-T0T/ito_web_app/internal/usecase"

	"github.com/gin-gonic/gin"
)

type Router struct {
	userUseCase *usecase.UserUseCase
	roomUseCase *usecase.RoomUseCase
	gameUseCase *usecase.GameUseCase
	logger      logger.ILogger
}

func NewRouter(userUC *usecase.UserUseCase, roomUC *usecase.RoomUseCase, gameUC *usecase.GameUseCase, log logger.ILogger) *Router {
	return &Router{
		userUseCase: userUC,
		roomUseCase: roomUC,
		gameUseCase: gameUC,
		logger:      log,
	}
}

func (r *Router) SetupRoutes(e *gin.Engine) {
	e.Use(middleware.LoggerMiddleware(r.logger))
	e.Use(middleware.CorsMiddleware())

	api := e.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", r.RegisterUser)
			auth.POST("/login", r.LoginUser)
		}

		secured := api.Group("")
		secured.Use(middleware.AuthMiddleware())
		{
			r.setupRoomRoutes(secured)
			r.setupGameRoutes(secured)
		}
	}
}

func (r *Router) setupRoomRoutes(rg *gin.RouterGroup) {
	rooms := rg.Group("/rooms")
	{
		rooms.GET("", r.GetRooms)
		rooms.POST("", r.validateCreateRoom, r.CreateRoom)
		rooms.GET("/:id", r.GetRoom)
		rooms.POST("/:id/join", r.validateJoinRoom, r.JoinRoom)
	}
}

func (r *Router) setupGameRoutes(rg *gin.RouterGroup) {
	games := rg.Group("/games")
	{
		games.POST("/:roomId/ready", r.SetPlayerReady)
		games.POST("/:roomId/start", r.StartGame)
		games.POST("/:roomId/vote", r.InitiateVote)
		games.GET("/:roomId/status", r.GetGameStatus)
		games.GET("/ws/:roomId", r.HandleWebSocket)
	}
}
