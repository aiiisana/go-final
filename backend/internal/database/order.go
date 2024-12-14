package database

import (
	"ecommerce-platform/internal/models"
	"errors"
	"log"

	"gorm.io/gorm"
)

func CreateOrder(db *gorm.DB, order *models.Order, currentUser *models.User) error {
	if currentUser.Role.RoleName == "guest" {
		return errors.New("guests cannot place orders")
	}

	order.UserID = currentUser.ID
	if err := db.Create(order).Error; err != nil {
		log.Println("Error creating order:", err)
		return err
	}
	return nil
}

func GetOrderByID(db *gorm.DB, orderID string) (*models.Order, error) {
	var order models.Order
	if err := db.Preload("OrderItems").Preload("Payments").First(&order, "id = ?", orderID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("order not found")
		}
		log.Println("Error fetching order:", err)
		return nil, err
	}
	return &order, nil
}

func GetOrdersByUser(db *gorm.DB, userID string) ([]models.Order, error) {
	var orders []models.Order
	if err := db.Preload("OrderItems").Preload("Payments").Where("user_id = ?", userID).Find(&orders).Error; err != nil {
		log.Println("Error fetching orders for user:", err)
		return nil, err
	}
	return orders, nil
}

func UpdateOrderStatus(db *gorm.DB, orderID string, status string, currentUser *models.User) error {
	if currentUser.Role.RoleName != "admin" {
		return errors.New("only admins can update order status")
	}

	if err := db.Model(&models.Order{}).Where("order_id = ?", orderID).Update("status", status).Error; err != nil {
		log.Println("Error updating order status:", err)
		return err
	}
	return nil
}

func DeleteOrder(db *gorm.DB, orderID string, currentUser *models.User) error {
	if currentUser.Role.RoleName != "admin" {
		return errors.New("only admins can delete orders")
	}
	if err := db.Delete(&models.Order{}, "order_id = ?", orderID).Error; err != nil {
		log.Println("Error deleting order:", err)
		return err
	}
	return nil
}

func GetProductPriceByID(db *gorm.DB, productID uint) (float64, error) {
	var product models.Product
	if err := db.First(&product, productID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, errors.New("product not found")
		}
		return 0, err
	}
	return product.Price, nil
}

func CreateOrderItem(db *gorm.DB, orderItem *models.OrderItem) error {
	if err := db.Create(orderItem).Error; err != nil {
		log.Println("Error creating order item:", err)
		return err
	}
	return nil
}

func UpdateOrderItem(db *gorm.DB, orderItemID uint, updatedData models.OrderItem) error {
	if err := db.Model(&models.OrderItem{}).Where("id = ?", orderItemID).Updates(updatedData).Error; err != nil {
		log.Println("Error updating order item:", err)
		return err
	}
	return nil
}

func GetOrderItems(db *gorm.DB, orderID uint) ([]models.OrderItem, error) {
	var orderItems []models.OrderItem
	if err := db.Where("order_id = ?", orderID).Find(&orderItems).Error; err != nil {
		log.Println("Error fetching order items:", err)
		return nil, err
	}
	return orderItems, nil
}

func CalculateOrderTotal(db *gorm.DB, orderID uint) (float64, error) {
	var totalAmount float64
	err := db.Model(&models.OrderItem{}).Where("order_id = ?", orderID).
		Select("COALESCE(SUM(price), 0)").Scan(&totalAmount).Error
	if err != nil {
		log.Println("Error calculating order total:", err)
		return 0, err
	}
	return totalAmount, nil
}

func UpdateOrderTotalAmount(db *gorm.DB, orderID uint, totalAmount float64) error {
	if err := db.Model(&models.Order{}).Where("id = ?", orderID).
		Update("total_amount", totalAmount).Error; err != nil {
		log.Println("Error updating order total amount:", err)
		return err
	}
	return nil
}
