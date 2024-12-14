package models

import (
	"time"

	"gorm.io/gorm"
)

type Review struct {
	gorm.Model
	ProductID uint      `json:"product_id"`
	UserID    uint      `json:"user_id"`
	Rating    float64   `json:"rating"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"created_at"`
}

type AuditLog struct {
	gorm.Model
	Action    string    `json:"action"`
	UserID    uint      `json:"user_id"`
	Timestamp time.Time `json:"timestamp"`
}
