package handlers

import (
	"main_videork/internal/application/useCase"
	"main_videork/internal/domain/interfaces"
	"main_videork/internal/presentation/middlewares"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewRouter(router *gin.Engine, authService *useCase.AuthService, userService *useCase.UserService, userRepo interfaces.UserRepository, secret string, uploadVideoUC *useCase.UploadVideoUseCase) {
	authHandlers := NewAuthHandlers(authService)
	userHandlers := NewUserHandlers(userService)
	videoHandlers := NewVideoHandlers(uploadVideoUC)

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	router.POST("/api/auth/signup", userHandlers.Register)
	router.POST("/api/auth/login", authHandlers.Login)

	authGroup := router.Group("/")
	authGroup.Use(middlewares.JWTMiddleware(authService, secret))
	authGroup.POST("/api/auth/logout", authHandlers.Logout)
	authGroup.GET("/api/me", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	videoGroup := authGroup.Group("/api/videos")
	videoGroup.Use(middlewares.PermissionMiddleware(userRepo, "upload_video"))
	videoGroup.POST("/upload", videoHandlers.Upload)

}
