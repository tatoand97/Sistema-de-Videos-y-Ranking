package http

import "github.com/gin-gonic/gin"

func SetupRouter(hUsers *UsersHandler, corsOrigin, jwtSecret string) *gin.Engine {
	r := gin.Default()
	r.Use(CORSMiddleware(corsOrigin))

	r.POST("/api/auth/signup", hUsers.Register)
	r.POST("/api/auth/login", hUsers.Login)

	auth := r.Group("/api", AuthRequired(jwtSecret))
	{
		auth.GET("/me", hUsers.Me)
	}
	return r
}
