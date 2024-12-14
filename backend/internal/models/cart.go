package models

import (
	"time"

	"gorm.io/gorm"
)

type Cart struct {
	gorm.Model
	UserID    uint       `json:"user_id"`
	CreatedAt time.Time  `json:"created_at"`
	CartItems []CartItem `gorm:"foreignKey:CartID;references:ID" json:"cart_items"`
}

type CartItem struct {
	gorm.Model
	CartID    uint `json:"cart_id"`
	ProductID uint `json:"product_id"`
	Quantity  int  `json:"quantity"`
}
