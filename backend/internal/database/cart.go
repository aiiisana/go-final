package database

import (
	"ecommerce-platform/internal/models"
	"errors"
	"log"

	"gorm.io/gorm"
)

func CreateCart(db *gorm.DB, cart *models.Cart, currentUser *models.User) error {
	cart.UserID = currentUser.ID
	if err := db.Create(cart).Error; err != nil {
		log.Println("Error creating cart:", err)
		return err
	}
	return nil
}

func GetCartByUserID(db *gorm.DB, userID string) (*models.Cart, error) {
	var cart models.Cart
	if err := db.Preload("CartItems").First(&cart, "id = ?", userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("cart not found")
		}
		log.Println("Error fetching cart:", err)
		return nil, err
	}
	return &cart, nil
}

func AddItemToCart(db *gorm.DB, item *models.CartItem) error {
	var cart models.Cart
	if err := db.First(&cart, "id = ?", item.CartID).Error; err != nil {
		return errors.New("cart not found")
	}

	var existingItem models.CartItem
	if err := db.Where("cart_id = ? AND product_id = ?", item.CartID, item.ProductID).First(&existingItem).Error; err == nil {
		existingItem.Quantity += item.Quantity
		if err := db.Save(&existingItem).Error; err != nil {
			log.Println("Error updating cart item:", err)
			return err
		}
		return nil
	}

	if err := db.Create(item).Error; err != nil {
		log.Println("Error adding item to cart:", err)
		return err
	}
	return nil
}

func DeleteCartItem(db *gorm.DB, itemID string) error {
	if err := db.Delete(&models.CartItem{}, "id = ?", itemID).Error; err != nil {
		log.Println("Error deleting cart item:", err)
		return err
	}
	return nil

}

func GetCartItemsByUserID(db *gorm.DB, userID string) ([]models.CartItem, error) {
	var cart models.Cart
	if err := db.Preload("CartItems").First(&cart, "user_id = ?", userID).Error; err != nil {
		log.Println("Error fetching cart:", err)
		return nil, err
	}

	return cart.CartItems, nil
}
