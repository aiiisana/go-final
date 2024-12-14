package handlers

import (
	"ecommerce-platform/internal/database"
	"ecommerce-platform/internal/models"
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
)

func CreateOrder(db database.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var order models.Order
		if err := c.ShouldBindJSON(&order); err != nil {
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

		if err := database.CreateOrder(db.DB, &order, user); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var wg sync.WaitGroup
		errCh := make(chan error, 2)

		wg.Add(1)
		go func() {
			defer wg.Done()
			totalAmount, err := database.CalculateOrderTotal(db.DB, order.ID)
			if err != nil {
				errCh <- err
				return
			}
			order.TotalAmount = totalAmount
			if err := database.UpdateOrderTotalAmount(db.DB, order.ID, totalAmount); err != nil {
				errCh <- err
			}
		}()

		wg.Wait()
		close(errCh)

		for err := range errCh {
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}

		c.JSON(http.StatusCreated, order)
	}
}

func CreateOrderItem(db database.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var orderItem models.OrderItem
		if err := c.ShouldBindJSON(&orderItem); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		var wg sync.WaitGroup
		errCh := make(chan error, 1)

		wg.Add(1)
		go func() {
			defer wg.Done()
			productPrice, err := database.GetProductPriceByID(db.DB, orderItem.ProductID)
			if err != nil {
				errCh <- err
				return
			}

			orderItem.Price = productPrice * float64(orderItem.Quantity)

			if err := database.CreateOrderItem(db.DB, &orderItem); err != nil {
				errCh <- err
			}
		}()

		wg.Wait()
		close(errCh)

		for err := range errCh {
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Order item created successfully", "order_item": orderItem})
	}
}

func UpdateOrderItem(db database.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		orderItemIDStr := c.Param("item_id")
		orderItemID, err := strconv.ParseUint(orderItemIDStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order item ID"})
			return
		}

		var updatedData models.OrderItem
		if err := c.ShouldBindJSON(&updatedData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		var wg sync.WaitGroup
		errCh := make(chan error, 1)

		wg.Add(1)
		go func() {
			defer wg.Done()
			productPrice, err := database.GetProductPriceByID(db.DB, updatedData.ProductID)
			if err != nil {
				errCh <- err
				return
			}

			updatedData.Price = productPrice * float64(updatedData.Quantity)

			if err := database.UpdateOrderItem(db.DB, uint(orderItemID), updatedData); err != nil {
				errCh <- err
			}
		}()

		wg.Wait()
		close(errCh)

		for err := range errCh {
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}

		c.JSON(http.StatusOK, gin.H{"message": "Order item updated successfully"})
	}
}

func GetOrderByID(db database.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		orderID := c.Param("id")

		order, err := database.GetOrderByID(db.DB, orderID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, order)
	}
}

func GetOrdersByUser(db database.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Param("user_id")

		orders, err := database.GetOrdersByUser(db.DB, userID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, orders)
	}
}

func GetOrderItems(db database.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		orderIDStr := c.Param("id")
		orderID, err := strconv.ParseUint(orderIDStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
			return
		}

		orderItems, err := database.GetOrderItems(db.DB, uint(orderID))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch order items"})
			return
		}

		totalAmount, err := database.CalculateOrderTotal(db.DB, uint(orderID))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to calculate total amount"})
			return
		}

		if err := database.UpdateOrderTotalAmount(db.DB, uint(orderID), totalAmount); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update total amount"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"order_items": orderItems, "total_amount": totalAmount})
	}
}
