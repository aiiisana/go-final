package database

import (
	"ecommerce-platform/internal/models"
	"errors"

	"gorm.io/gorm"
)

func CreateCategory(db *gorm.DB, category *models.Category) error {
	if err := db.Create(category).Error; err != nil {
		return err
	}
	return nil
}

func GetCategories(db *gorm.DB) ([]models.Category, error) {
	var categories []models.Category
	if err := db.Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

func GetCategoryByID(db *gorm.DB, categoryID uint) (*models.Category, error) {
	var category models.Category
	if err := db.First(&category, categoryID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &category, nil
}

func UpdateCategory(db *gorm.DB, categoryID uint, updatedCategory models.Category) error {
	if err := db.Model(&models.Category{}).Where("id = ?", categoryID).Updates(updatedCategory).Error; err != nil {
		return err
	}
	return nil
}

func DeleteCategory(db *gorm.DB, categoryID uint) error {
	if err := db.Delete(&models.Category{}, "id = ?", categoryID).Error; err != nil {
		return err
	}
	return nil
}
