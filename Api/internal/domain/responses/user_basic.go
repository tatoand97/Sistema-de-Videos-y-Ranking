package responses

// UserBasic contiene informacion minima para ranking: username y ciudad.
type UserBasic struct {
	UserID   uint    `json:"-" gorm:"column:user_id"`
	Username string  `json:"username" gorm:"column:username"`
	City     *string `json:"city,omitempty" gorm:"column:city"`
}
