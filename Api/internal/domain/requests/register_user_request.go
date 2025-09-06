package requests

type RegisterUserRequest struct {
	FirstName string `json:"first_name" form:"first_name" binding:"required"`
	LastName  string `json:"last_name" form:"last_name" binding:"required"`
	Email     string `json:"email" form:"email" binding:"required,email"`
	Password1 string `json:"password1" form:"password1" binding:"required"`
	Password2 string `json:"password2" form:"password2" binding:"required"`
	City      string `json:"city" form:"city" binding:"required"`
	Country   string `json:"country" form:"country" binding:"required"`
}
