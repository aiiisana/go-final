package handlers

import (
	"ecommerce-platform/internal/database"
	"ecommerce-platform/internal/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetUserSessions(db database.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDStr := c.Param("user_id")
		userID, err := strconv.ParseUint(userIDStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		sessions, err := database.GetUserSessions(db.DB, uint(userID))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching sessions"})
			return
		}

		validSessions := []models.Session{}
		for _, session := range sessions {
			if !session.IsExpired() {
				validSessions = append(validSessions, session)
			}
		}

		c.JSON(http.StatusOK, validSessions)
	}
}

func GetAuditLogs(db database.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDStr := c.Param("user_id")
		userID, err := strconv.ParseUint(userIDStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		logs, err := database.GetAuditLogs(db.DB, uint(userID))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching audit logs"})
			return
		}

		c.JSON(http.StatusOK, logs)
	}
}
