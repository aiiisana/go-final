package handlers

import (
	"ecommerce-platform/internal/database"
	"ecommerce-platform/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateReview(db database.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var review models.Review
		if err := c.ShouldBindJSON(&review); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		currentUser, exists := c.Get("current_user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		user, ok := currentUser.(*models.User)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user data in context"})
			return
		}

		review.UserID = user.ID
		if err := database.CreateReview(db.DB, &review); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating review"})
			return
		}

		c.JSON(http.StatusCreated, review)
	}
}

func GetReviewsByProduct(db database.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		productID := c.Param("product_id")

		reviews, err := database.GetReviewsByProduct(db.DB, productID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "No reviews found"})
			return
		}

		c.JSON(http.StatusOK, reviews)
	}
}

func DeleteReview(db database.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		reviewID := c.Param("id")

		currentUser, exists := c.Get("current_user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		user, ok := currentUser.(*models.User)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user data in context"})
			return
		}

		if err := database.DeleteReview(db.DB, reviewID, user); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting review"})
			return
		}

		c.JSON(http.StatusNoContent, nil)
	}
}
