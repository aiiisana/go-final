package models

import (
	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	UserID      uint        `json:"user_id"`
	Status      string      `json:"status"`
	TotalAmount float64     `json:"total_amount"`
	OrderItems  []OrderItem `gorm:"foreignKey:OrderID"`
	Payments    []Payment   `gorm:"foreignKey:OrderID"`
}

type OrderItem struct {
	gorm.Model
	OrderID   uint    `json:"order_id"`
	ProductID uint    `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

type Payment struct {
	gorm.Model
	OrderID       uint    `json:"order_id"`
	Amount        float64 `json:"amount"`
	PaymentMethod string  `json:"payment_method"`
}
