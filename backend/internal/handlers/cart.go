package handlers

import (
	"ecommerce-platform/internal/database"
	"ecommerce-platform/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateShoppingCart(db database.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var cart models.Cart
		if err := c.ShouldBindJSON(&cart); err != nil {
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

		if err := database.CreateCart(db.DB, &cart, user); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, cart)
	}
}

func GetCartByUserID(db database.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Param("user_id")

		cart, err := database.GetCartByUserID(db.DB, userID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, cart)
	}
}

func AddItemToCart(db database.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var item models.CartItem
		if err := c.ShouldBindJSON(&item); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		if err := database.AddItemToCart(db.DB, &item); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, item)
	}
}

func DeleteCartItem(db database.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		itemID := c.Param("item_id")

		if err := database.DeleteCartItem(db.DB, itemID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusNoContent, nil)
	}
}

func GetCartItemsByUserID(db database.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Param("user_id")

		if userID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
			return
		}

		items, err := database.GetCartItemsByUserID(db.DB, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, items)
	}
}
