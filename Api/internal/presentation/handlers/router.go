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
}

func NewRouter(router *gin.Engine, cfg RouterConfig) {
	authHandlers := NewAuthHandlers(cfg.AuthService)
	userHandlers := NewUserHandlers(cfg.UserService)
	videoHandlers := NewVideoHandlers(cfg.UploadsUC)
	locationHandlers := NewLocationHandlers(cfg.LocationService)
	// Constructor con cache de solo lectura
	publicHandlers := NewPublicHandlersWithCache(cfg.PublicService, cfg.Cache)

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	// Publico: ubicacion
	router.GET("/api/location/city-id", locationHandlers.GetCityID)
	router.GET("/api/public/videos", publicHandlers.ListPublicVideos)
	router.GET("/api/public/rankings", publicHandlers.ListRankings)
	// Se eliminaron endpoints basados en poll_id (leaderboard/stats/count)
	router.POST("/api/auth/signup", userHandlers.Register)
	router.POST("/api/auth/login", authHandlers.Login)

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

	// Ruta protegida para votar por un video publico
	authGroup.POST("/api/public/videos/:video_id/vote", publicHandlers.VotePublicVideo)

}
