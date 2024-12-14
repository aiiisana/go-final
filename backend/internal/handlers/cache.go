package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"ecommerce-platform/internal/database"
	"ecommerce-platform/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetProductFromCache(client database.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		cacheKey := c.Param("cache_key")

		cachedData, err := database.GetCacheRedis(cacheKey)
		if err == nil && cachedData != "" {
			var product models.Product
			err = json.Unmarshal([]byte(cachedData), &product)
			if err != nil {
				log.Println("Error unmarshalling cached product data:", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error parsing cached data"})
				return
			}
			c.JSON(http.StatusOK, product)
			return
		}

		var product models.Product
		if err := client.DB.Where("id = ?", cacheKey).First(&product).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "Product not found"})
			return
		}

		productData, err := json.Marshal(product)
		if err != nil {
			log.Println("Error marshalling product:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error marshalling product"})
			return
		}

		err = database.SetCacheRedis(cacheKey, string(productData), time.Hour)
		if err != nil {
			log.Println("Error setting product data in Redis cache:", err)
		}

		c.JSON(http.StatusOK, product)
	}
}

func DeleteCache(client *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		cacheKey := c.Param("key")

		err := database.DeleteCacheRedis(cacheKey)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete cache"})
			return
		}

		var cacheEntry models.Cache
		if err := client.Where("cache_key = ?", cacheKey).Delete(&cacheEntry).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete cache from database"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Cache deleted successfully"})
	}
}
