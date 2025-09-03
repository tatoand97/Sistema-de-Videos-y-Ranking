package requests

type RegisterUserRequest struct {
	FirstName string `json:"first_name" form:"first_name" binding:"required"`
	LastName  string `json:"last_name" form:"last_name" binding:"required"`
	Email     string `json:"email" form:"email" binding:"required,email"`
	Password1 string `json:"password1" form:"password1" binding:"required"`
	Password2 string `json:"password2" form:"password2" binding:"required"`
	CityID    int    `json:"city_id" form:"city_id" binding:"required"`
}
