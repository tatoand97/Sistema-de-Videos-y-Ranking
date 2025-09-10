package handlers

import (
	"api/internal/application/useCase"
	"api/internal/domain/interfaces"
	"api/internal/presentation/middlewares"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RouterConfig groups router dependencies to reduce parameters.
type RouterConfig struct {
	AuthService     *useCase.AuthService
	UserService     *useCase.UserService
	LocationService *useCase.LocationService
	UploadsUC       *useCase.UploadsUseCase
	PublicService   *useCase.PublicService
	StatusService   *useCase.StatusService
	JWTSecret       string
	Cache           interfaces.Cache
	IdemTTLSeconds  int
	Aggregates      interfaces.Aggregates
}

func NewRouter(router *gin.Engine, cfg RouterConfig) {
	authHandlers := NewAuthHandlers(cfg.AuthService)
	userHandlers := NewUserHandlers(cfg.UserService)
	videoHandlers := NewVideoHandlers(cfg.UploadsUC)
	locationHandlers := NewLocationHandlers(cfg.LocationService)
	// Prefer constructor with cache & aggregates; keep simple one available for tests
	publicHandlers := NewPublicHandlersFull(cfg.PublicService, cfg.Cache, cfg.IdemTTLSeconds, cfg.Aggregates)
	statusHandlers := NewStatusHandlers(cfg.StatusService)

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	// Público: ubicación
	router.GET("/api/location/city-id", locationHandlers.GetCityID)
	router.GET("/api/public/videos", publicHandlers.ListPublicVideos)
	router.GET("/api/public/rankings", publicHandlers.ListRankings)
	// Leaderboards/Stats via Redis aggregates
	router.GET("/api/public/leaderboard/:poll_id", publicHandlers.GetLeaderboard)
	router.GET("/api/public/stats/:poll_id", publicHandlers.GetStats)
	router.GET("/api/public/count/:poll_id/:member", publicHandlers.GetCount)
	router.POST("/api/auth/signup", userHandlers.Register)
	router.POST("/api/auth/login", authHandlers.Login)

	// Public: list video statuses
	router.GET("/api/videos/statuses", statusHandlers.ListVideoStatuses)

	authGroup := router.Group("/")
	authGroup.Use(middlewares.JWTMiddleware(cfg.AuthService, cfg.JWTSecret))
	authGroup.POST("/api/auth/logout", authHandlers.Logout)
	authGroup.GET("/api/me", authHandlers.Me)

	videoGroup := authGroup.Group("/api/videos")
	videoGroup.GET("", videoHandlers.ListVideos)
	videoGroup.POST("/upload", videoHandlers.Upload)
	videoGroup.GET("/:video_id", videoHandlers.GetVideoDetail)
	videoGroup.DELETE("/:video_id", videoHandlers.DeleteVideo)
	videoGroup.POST("/:video_id/publish", videoHandlers.PublishVideo)

	// Ruta protegida para votar por un video público
	authGroup.POST("/api/public/videos/:video_id/vote", publicHandlers.VotePublicVideo)

}
