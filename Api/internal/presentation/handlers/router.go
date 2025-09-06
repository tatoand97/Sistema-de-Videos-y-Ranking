package handlers

import (
	"main_videork/internal/application/useCase"
	"main_videork/internal/presentation/middlewares"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewRouter(router *gin.Engine, authService *useCase.AuthService, userService *useCase.UserService, locationService *useCase.LocationService, secret string, uploadVideoUC *useCase.UploadVideoUseCase) {
	authHandlers := NewAuthHandlers(authService)
	userHandlers := NewUserHandlers(userService)
	videoHandlers := NewVideoHandlers(uploadVideoUC)
	locationHandlers := NewLocationHandlers(locationService)

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	// Público: ubicación
	router.GET("/api/location/city-id", locationHandlers.GetCityID)
	router.POST("/api/auth/signup", userHandlers.Register)
	router.POST("/api/auth/login", authHandlers.Login)

	authGroup := router.Group("/")
	authGroup.Use(middlewares.JWTMiddleware(authService, secret))
	authGroup.POST("/api/auth/logout", authHandlers.Logout)
	authGroup.GET("/api/me", authHandlers.Me)

	videoGroup := authGroup.Group("/api/videos")
	videoGroup.POST("/upload", videoHandlers.Upload)

}
