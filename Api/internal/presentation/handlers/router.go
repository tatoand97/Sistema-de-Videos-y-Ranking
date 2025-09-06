package handlers

import (
	"api/internal/application/useCase"
	"api/internal/presentation/middlewares"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewRouter(router *gin.Engine, authService *useCase.AuthService, userService *useCase.UserService, locationService *useCase.LocationService, secret string, uploadVideoUC *useCase.UploadVideoUseCase, publicService *useCase.PublicService, statusService *useCase.StatusService) {
	authHandlers := NewAuthHandlers(authService)
	userHandlers := NewUserHandlers(userService)
	videoHandlers := NewVideoHandlers(uploadVideoUC)
	locationHandlers := NewLocationHandlers(locationService)
	publicHandlers := NewPublicHandlers(publicService)
	statusHandlers := NewStatusHandlers(statusService)

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	// Público: ubicación
	router.GET("/api/location/city-id", locationHandlers.GetCityID)
	router.GET("/api/public/videos", publicHandlers.ListPublicVideos)
	router.GET("/api/public/rankings", publicHandlers.ListRankings)
	router.POST("/api/auth/signup", userHandlers.Register)
	router.POST("/api/auth/login", authHandlers.Login)

	// Public: list video statuses
	router.GET("/api/videos/statuses", statusHandlers.ListVideoStatuses)

	authGroup := router.Group("/")
	authGroup.Use(middlewares.JWTMiddleware(authService, secret))
	authGroup.POST("/api/auth/logout", authHandlers.Logout)
	authGroup.GET("/api/me", authHandlers.Me)

	videoGroup := authGroup.Group("/api/videos")
	videoGroup.POST("/upload", videoHandlers.Upload)

	// Ruta protegida para votar por un video público
	authGroup.POST("/api/public/videos/:video_id/vote", publicHandlers.VotePublicVideo)

}
