package entities

import "time"

type User struct {
	UserID       int       `gorm:"primaryKey;column:user_id"`
	FirstName    string    `gorm:"column:first_name"`
	LastName     string    `gorm:"column:last_name"`
	Email        string    `gorm:"column:email;uniqueIndex;not null"`
	PasswordHash string    `gorm:"column:password_hash;not null"`
	CityID       int       `gorm:"column:city_id;not null"`
	CreatedAt    time.Time `gorm:"column:created_at"`
}

func (User) TableName() string {
	return "users"
}
