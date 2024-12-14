package handlers

import (
	"ecommerce-platform/internal/auth"
	"ecommerce-platform/internal/database"
	logging "ecommerce-platform/internal/logger"
	"ecommerce-platform/internal/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Login(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var loginRequest struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		if err := c.ShouldBindJSON(&loginRequest); err != nil {
			logging.LogError("Invalid login request", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		var user models.User
		if err := db.Where("username = ?", loginRequest.Username).First(&user).Error; err != nil {
			logging.LogError("User not found", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
			return
		}

		if !CheckPasswordHash(loginRequest.Password, user.Password) {
			logging.LogError("Invalid password", nil)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
			return
		}

		role := user.Role.RoleName
		token, err := auth.GenerateToken(user.ID, role)
		if err != nil {
			logging.LogError("Error generating token", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
			return
		}

		sessionExpiresAt := time.Now().Add(24 * time.Hour)
		session, err := database.CreateSession(db, user.ID, sessionExpiresAt)
		if err != nil {
			logging.LogError("Error creating session", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create session"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"token": token,
			"session": gin.H{
				"id":         session.ID,
				"user_id":    session.UserID,
				"expires_at": session.ExpiresAt,
			},
		})
	}
}

func Logout(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionID := c.Param("session_id")

		var session models.Session
		if err := db.Where("id = ?", sessionID).First(&session).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid session"})
			logging.LogError("Failed logout attempt - Invalid session ID", err)
			return
		}

		if err := db.Delete(&session).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting session"})
			logging.LogError("Failed to delete session", err)
			return
		}

		logging.LogInfo("User logged out successfully: " + sessionID)

		c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
	}
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
