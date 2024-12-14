// internal/database/product.go
package database

import (
	"ecommerce-platform/internal/models"
	"errors"
	"log"

	"gorm.io/gorm"
)

func CreateProduct(db *gorm.DB, product *models.Product, currentUser *models.User) error {
	if currentUser.Role.RoleName != "admin" {
		return errors.New("only admins can create products")
	}

	if err := db.Create(product).Error; err != nil {
		log.Println("Error creating product:", err)
		return err
	}
	return nil
}

func UpdateProduct(db *gorm.DB, productID string, updatedData models.Product) error {
	if err := db.Model(&models.Product{}).Where("id = ?", productID).Updates(updatedData).Error; err != nil {
		log.Println("Error updating product:", err)
		return err
	}
	return nil
}

func DeleteProduct(db *gorm.DB, productID string) error {

	if err := db.Delete(&models.Product{}, "id = ?", productID).Error; err != nil {
		log.Println("Error deleting product:", err)
		return err
	}
	return nil
}

func GetProducts(db *gorm.DB) ([]models.Product, error) {
	var products []models.Product
	if err := db.Find(&products).Error; err != nil {
		log.Println("Error fetching products:", err)
		return nil, err
	}
	return products, nil
}

func GetProductByID(db *gorm.DB, productID string) (*models.Product, error) {
	var product models.Product
	if err := db.First(&product, productID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &product, nil
}

func GetProductsByCategory(db *gorm.DB, categoryID string) ([]models.Product, error) {
	var products []models.Product
	if err := db.Preload("Category").Where("category_id = ?", categoryID).Find(&products).Error; err != nil {
		log.Println("Error fetching products by category:", err)
		return nil, err
	}
	return products, nil
}

func AddProductImage(db *gorm.DB, productID uint, imageURL string) (*models.ProductImage, error) {
	productImage := models.ProductImage{
		ProductID: productID,
		ImageURL:  imageURL,
	}

	if err := db.Create(&productImage).Error; err != nil {
		return nil, err
	}

	return &productImage, nil
}

func GetProductImages(db *gorm.DB, productID string) ([]models.ProductImage, error) {
	var images []models.ProductImage
	if err := db.Where("product_id = ?", productID).Find(&images).Error; err != nil {
		return nil, err
	}
	return images, nil
}

func UpdateProductImage(db *gorm.DB, imageID string, imageURL string) (*models.ProductImage, error) {
	var productImage models.ProductImage
	if err := db.Where("id = ?", imageID).First(&productImage).Error; err != nil {
		return nil, err
	}

	productImage.ImageURL = imageURL
	if err := db.Save(&productImage).Error; err != nil {
		return nil, err
	}

	return &productImage, nil
}

func DeleteProductImage(db *gorm.DB, imageID string) error {
	if err := db.Where("id = ?", imageID).Delete(&models.ProductImage{}).Error; err != nil {
		return err
	}
	return nil
}
