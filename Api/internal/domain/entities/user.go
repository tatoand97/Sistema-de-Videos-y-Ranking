package entities

import "time"

type User struct {
	UserId       uint   `gorm:"primaryKey"`
	FirstName    string `gorm:"not null"`
	LastName     string `gorm:"not null"`
	Email        string `gorm:"uniqueIndex;not null"`
	PasswordHash string `gorm:"not null"`
	City         string
	Country      string
	CreatedAt    time.Time
}
