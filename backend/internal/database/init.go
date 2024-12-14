package database

import (
	"ecommerce-platform/internal/models"
	"log"
)

func InitializeDatabase() (*Client, error) {
	client, err := NewDBClient()
	if err != nil {
		log.Fatal("Error connecting to the database:", err)
		return nil, err
	}

	err = client.RunMigration()
	if err != nil {
		log.Fatal("Error running migrations:", err)
		return nil, err
	}

	CreateAdminRole(client)
	CreateUserRole(client)

	importUsers(client, "./data/users.json")
	importData(client, "./data/categories.json", []models.Category{})

	importData(client, "./data/products.json", []models.Product{})
	importData(client, "./data/orders.json", []models.Order{})
	importData(client, "./data/carts.json", []models.Cart{})
	importData(client, "./data/payments.json", []models.Payment{})
	importData(client, "./data/order_item.json", []models.OrderItem{})
	importData(client, "./data/cart_items.json", []models.CartItem{})
	importData(client, "./data/reviews.json", []models.Review{})
	importData(client, "./data/user_addresses.json", []models.UserAddress{})
	importData(client, "./data/product_images.json", []models.ProductImage{})

	return client, nil
}

func CreateAdminRole(db *Client) {
	var role models.Role
	if err := db.DB.Where("role_name = ?", "admin").First(&role).Error; err != nil {
		role = models.Role{
			RoleName: "admin",
		}
		role.ID = 1
		if err := db.DB.Create(&role).Error; err != nil {
			log.Fatal("Failed to create admin role:", err)
		}
	}
}

func CreateUserRole(db *Client) {
	var role models.Role
	if err := db.DB.Where("role_name = ?", "user").First(&role).Error; err != nil {
		role = models.Role{
			RoleName: "user",
		}
		role.ID = 2
		if err := db.DB.Create(&role).Error; err != nil {
			log.Fatal("Failed to create user role:", err)
		}
	}
}
