package handlers

import (
	"ecommerce-platform/internal/database"
	"ecommerce-platform/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreatePayment(db database.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var payment models.Payment
		if err := c.ShouldBindJSON(&payment); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		if err := database.CreatePayment(db.DB, &payment); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating payment"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Payment created successfully", "payment": payment})
	}
}

func GetPayments(db database.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		payments, err := database.GetPayments(db.DB)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching payments"})
			return
		}
		c.JSON(http.StatusOK, payments)
	}
}

func GetPaymentByID(db database.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		paymentID := c.Param("id")

		payment, err := database.GetPaymentByID(db.DB, paymentID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching payment"})
			return
		}

		if payment == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Payment not found"})
			return
		}

		c.JSON(http.StatusOK, payment)
	}
}

func DeletePayment(db database.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		paymentID := c.Param("id")

		if err := database.DeletePayment(db.DB, paymentID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting payment"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Payment deleted successfully"})
	}
}
