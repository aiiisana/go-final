package handlers

import (
	"ecommerce-platform/internal/database"
	"ecommerce-platform/internal/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateProduct(db database.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var product models.Product

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

		if err := c.ShouldBindJSON(&product); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product data"})
			return
		}

		if err := database.CreateProduct(db.DB, &product, user); err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Product created successfully", "product": product})
	}
}

func GetAllProducts(db database.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		products, err := database.GetProducts(db.DB)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching products"})
			return
		}
		c.JSON(http.StatusOK, products)
	}
}

func GetProductByID(db database.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		productID := c.Param("product_id")
		product, err := database.GetProductByID(db.DB, productID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			return
		}
		c.JSON(http.StatusOK, product)
	}
}

func GetProductsByCategory(db database.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		categoryID := c.Param("category_id")

		products, err := database.GetProductsByCategory(db.DB, categoryID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching products by category"})
			return
		}

		c.JSON(http.StatusOK, products)
	}
}

func UpdateProduct(db database.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		productID := c.Param("id")

		var updatedProduct models.Product
		if err := c.ShouldBindJSON(&updatedProduct); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category data"})
			return
		}

		if err := database.UpdateProduct(db.DB, productID, updatedProduct); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update category"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Product updated successfully"})
	}
}

func DeleteProduct(db database.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		productID := c.Param("id")

		if err := database.DeleteProduct(db.DB, productID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete product"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
	}
}

func AddProductImage(db database.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		productIDStr := c.Param("product_id")
		productID, err := strconv.ParseUint(productIDStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
			return
		}

		var input struct {
			ImageURL string `json:"image_url"`
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		productImage, err := database.AddProductImage(db.DB, uint(productID), input.ImageURL)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add image"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Image added successfully", "image": productImage})
	}
}

func GetProductImages(db database.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		productID := c.Param("product_id")

		images, err := database.GetProductImages(db.DB, productID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch images"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"images": images})
	}
}

func UpdateProductImage(db database.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		imageID := c.Param("id")
		var input struct {
			ImageURL string `json:"image_url"`
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		updatedImage, err := database.UpdateProductImage(db.DB, imageID, input.ImageURL)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update image"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Image updated successfully", "image": updatedImage})
	}
}

func DeleteProductImage(db database.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		imageID := c.Param("id")

		if err := database.DeleteProductImage(db.DB, imageID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete image"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Image deleted successfully"})
	}
}
