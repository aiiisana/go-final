package database

import (
	"ecommerce-platform/internal/models"
	"errors"
	"log"

	"gorm.io/gorm"
)

func CreateReview(db *gorm.DB, review *models.Review) error {
	if err := db.Create(review).Error; err != nil {
		log.Println("Error creating review:", err)
		return err
	}
	return nil
}

func GetReviewsByProduct(db *gorm.DB, productID string) ([]models.Review, error) {
	var reviews []models.Review
	if err := db.Where("product_id = ?", productID).Find(&reviews).Error; err != nil {
		log.Println("Error fetching reviews for product:", err)
		return nil, err
	}
	return reviews, nil
}

func DeleteReview(db *gorm.DB, reviewID string, currentUser *models.User) error {
	var review models.Review
	if err := db.First(&review, "review_id = ?", reviewID).Error; err != nil {
		log.Println("Error fetching review:", err)
		return err
	}

	if review.UserID != currentUser.ID && currentUser.Role.RoleName != "admin" {
		return errors.New("you can only delete your own review or if you are an admin")
	}

	if err := db.Delete(&models.Review{}, "review_id = ?", reviewID).Error; err != nil {
		log.Println("Error deleting review:", err)
		return err
	}
	return nil
}
