package http

import (
	"main_prj/internal/ports"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

type UsersHandler struct {
	Auth      ports.AuthService
	UploadDir string
}

func (h *UsersHandler) Register(c *gin.Context) {
	username := c.PostForm("nombreUsuario")
	email := c.PostForm("email")
	password := c.PostForm("contrasena")

	var pathPtr *string
	if file, err := c.FormFile("imagenPerfil"); err == nil && file != nil {
		dst := filepath.Join(h.UploadDir, file.Filename)
		if err := c.SaveUploadedFile(file, dst); err == nil {
			path := dst
			pathPtr = &path
		}
	}

	user, err := h.Auth.Register(username, email, password, pathPtr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, user)
}

func (h *UsersHandler) Login(c *gin.Context) {
	var req struct {
		NombreUsuario string `json:"Nombre de usuario" form:"Nombre de usuario"`
		Contrasena    string `json:"Contraseña" form:"Contraseña"`
		Username      string `json:"username"`
		Password      string `json:"password"`
		Email         string `json:"email"`
	}
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}
	// username := req.Username
	// if username == "" {
	// 	username = req.NombreUsuario
	// }
	email := req.Email
	if email == "" {
		email = req.Email
	}
	password := req.Password
	if password == "" {
		password = req.Contrasena
	}

	token, user, err := h.Auth.Login(email, password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token, "usuario": user})
}

func (h *UsersHandler) Me(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"user_id": c.GetInt("user_id"),
		"email":   c.GetString("email"),
	})
}
