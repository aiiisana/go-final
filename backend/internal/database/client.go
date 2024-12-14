package database

import (
	"ecommerce-platform/internal/models"
	"fmt"
	"log"
	"os"
	"strconv"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Client struct {
	DB *gorm.DB
}

func NewDBClient() (*Client, error) {
	dbHost := getEnv("DB_HOST", "localhost")
	dbUsername := getEnv("DB_USERNAME", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "superuserPostgre")
	dbName := getEnv("DB_NAME", "finalgo")
	dbPort := getEnv("DB_PORT", "5432")

	databasePort, err := strconv.Atoi(dbPort)
	if err != nil {
		log.Fatal("Invalid DB Port:", err)
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		dbHost, dbUsername, dbPassword, dbName, databasePort, "disable")

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	client := Client{DB: db}
	return &client, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func (c *Client) RunMigration() error {

	if err := c.DB.AutoMigrate(
		&models.Role{},
	); err != nil {
		return fmt.Errorf("step 1 migration failed: %w", err)
	}
	if err := c.DB.AutoMigrate(
		&models.User{},
	); err != nil {
		return fmt.Errorf("step 2 migration failed: %w", err)
	}

	if err := c.DB.AutoMigrate(
		&models.Category{},
	); err != nil {
		log.Printf("Migration for Category failed: %v", err)
		return fmt.Errorf("step 3 migration failed: %w", err)
	}

	if err := c.DB.AutoMigrate(
		&models.Product{},
	); err != nil {
		return fmt.Errorf("step 4 migration failed: %w", err)
	}

	if err := c.DB.AutoMigrate(
		&models.ProductImage{},
	); err != nil {
		return fmt.Errorf("step 5 migration failed: %w", err)
	}

	if err := c.DB.AutoMigrate(
		&models.Order{},
	); err != nil {
		return fmt.Errorf("step 6 migration failed: %w", err)
	}

	if err := c.DB.AutoMigrate(
		&models.Cart{},
	); err != nil {
		return fmt.Errorf("step 7 migration failed: %w", err)
	}
	if err := c.DB.AutoMigrate(
		&models.OrderItem{},
		&models.Payment{},
		&models.CartItem{},
		&models.UserAddress{},
		&models.Review{},
		&models.Session{},
		&models.AuditLog{},
		&models.Cache{},
	); err != nil {
		return fmt.Errorf("step 8 migration failed: %w", err)
	}

	return nil
}

func (c *Client) CloseConnection() {
	sqlDB, _ := c.DB.DB()
	err := sqlDB.Close()
	if err != nil {
		log.Fatal("Failed to close database connection:", err)
	}
}
