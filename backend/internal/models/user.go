package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username  string    `gorm:"uniqueIndex;size:255" json:"username"`
	Password  string    `json:"password"`
	Email     string    `gorm:"uniqueIndex;size:255" json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	RoleID    uint      `json:"role_id"`
	Role      Role      `gorm:"foreignKey:RoleID" json:"role"`
}

type Role struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	RoleName string `gorm:"uniqueIndex;size:255" json:"role_name"`
}

type UserAddress struct {
	gorm.Model
	UserID  uint   `json:"user_id"`
	Street  string `json:"street"`
	City    string `json:"city"`
	State   string `json:"state"`
	ZipCode string `json:"zip_code"`
}

type Session struct {
	gorm.Model
	UserID    uint      `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

func (s *Session) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}
