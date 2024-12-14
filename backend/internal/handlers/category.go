package handlers

import (
	"ecommerce-platform/internal/database"
	"ecommerce-platform/internal/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateCategory(db database.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var category models.Category

		if err := c.ShouldBindJSON(&category); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category data"})
			return
		}

		if err := database.CreateCategory(db.DB, &category); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create category"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Category created successfully", "category": category})
	}
}

func GetAllCategories(db database.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		categories, err := database.GetCategories(db.DB)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch categories"})
			return
		}

		c.JSON(http.StatusOK, categories)
	}
}

func GetCategoryByID(db database.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		categoryIDStr := c.Param("id")
		categoryID, err := strconv.ParseUint(categoryIDStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
			return
		}

		category, err := database.GetCategoryByID(db.DB, uint(categoryID))
		if err != nil || category == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
			return
		}

		c.JSON(http.StatusOK, category)
	}
}

func UpdateCategory(db database.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		categoryIDStr := c.Param("id")
		categoryID, err := strconv.ParseUint(categoryIDStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
			return
		}

		var updatedCategory models.Category
		if err := c.ShouldBindJSON(&updatedCategory); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category data"})
			return
		}

		if err := database.UpdateCategory(db.DB, uint(categoryID), updatedCategory); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update category"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Category updated successfully"})
	}
}

func DeleteCategory(db database.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		categoryIDStr := c.Param("id")
		categoryID, err := strconv.ParseUint(categoryIDStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
			return
		}

		if err := database.DeleteCategory(db.DB, uint(categoryID)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete category"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Category deleted successfully"})
	}
}
