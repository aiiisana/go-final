package handlers

import (
	"ecommerce-platform/internal/database"
	"ecommerce-platform/internal/models"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(db database.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid input",
				"details": err.Error(),
			})
			return
		}

		if len(user.Password) < 6 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Password must be at least 6 characters"})
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
			return
		}

		user.Password = string(hashedPassword)

		var userRole models.Role
		if err := db.DB.Where("role_name = ?", "user").First(&userRole).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "User role not found"})
			return
		}
		user.RoleID = userRole.ID

		if err := database.CreateUser(db.DB, &user); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
	}
}

func CreateAdminUser(db database.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid input",
				"details": err.Error(),
			})
			return
		}

		if len(user.Password) < 6 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Password must be at least 6 characters"})
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
			return
		}

		user.Password = string(hashedPassword)

		var adminRole models.Role
		if err := db.DB.Where("role_name = ?", "admin").First(&adminRole).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Admin role not found"})
			return
		}
		user.RoleID = adminRole.ID

		if err := database.CreateUser(db.DB, &user); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Admin user created successfully"})
	}
}

func CreateUserForAdmin(db database.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		currentUser, exists := c.Get("current_user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		user, ok := currentUser.(*models.User)
		if !ok || user.Role.RoleName != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to create a user"})
			return
		}

		var newUser models.User
		if err := c.ShouldBindJSON(&newUser); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		if len(newUser.Password) < 6 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Password must be at least 6 characters"})
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
			return
		}

		newUser.Password = string(hashedPassword)

		if err := database.CreateUser(db.DB, &newUser); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
	}
}

func GetUserByID(client database.Client, userID uint) (*models.User, error) {
	var user models.User
	if err := client.DB.Preload("Role").Where("id = ?", userID).First(&user).Error; err != nil {
		log.Println("Error fetching user:", err)
		return nil, err
	}

	return &user, nil
}

func UpdateUser(db database.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDStr := c.Param("id")

		userID, err := strconv.ParseUint(userIDStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		var updatedData models.User
		if err := c.ShouldBindJSON(&updatedData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		currentUser := c.MustGet("current_user").(*models.User)

		if err := database.UpdateUser(db.DB, uint(userID), updatedData, currentUser); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
	}
}

func DeleteUser(db database.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDStr := c.Param("id")

		userID, err := strconv.ParseUint(userIDStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		currentUser := c.MustGet("current_user").(*models.User)

		if err := database.DeleteUser(db.DB, uint(userID), currentUser); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
	}
}

func GetAllUsers(db database.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		currentUser, exists := c.Get("current_user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		user, ok := currentUser.(*models.User)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user data"})
			return
		}

		users, err := database.GetAllUsers(db.DB, user)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, users)
	}
}

func GetProfile(c *gin.Context) {
	currentUser, exists := c.Get("current_user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	user, ok := currentUser.(*models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id":    user.ID,
		"username":   user.Username,
		"email":      user.Email,
		"role":       user.Role.RoleName,
		"created_at": user.CreatedAt,
	})
}

func CreateUserAddress(db database.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var address models.UserAddress
		if err := c.ShouldBindJSON(&address); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		createdAddress, err := database.CreateUserAddress(db.DB, address.UserID, address.Street, address.City, address.State, address.ZipCode)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating address"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"message": "Address created successfully",
			"address": createdAddress,
		})
	}
}

func GetUserAddresses(db database.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
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

		addresses, err := database.GetUserAddresses(db.DB, user.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching addresses"})
			return
		}

		c.JSON(http.StatusOK, addresses)
	}
}

func UpdateUserAddress(db database.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		addressIDStr := c.Param("id")
		addressID, err := strconv.ParseUint(addressIDStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid address ID"})
			return
		}

		var updatedAddress models.UserAddress
		if err := c.ShouldBindJSON(&updatedAddress); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		if err := database.UpdateUserAddress(db.DB, uint(addressID), updatedAddress); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating address"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Address updated successfully"})
	}
}

func DeleteUserAddress(db database.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		addressIDStr := c.Param("id")
		addressID, err := strconv.ParseUint(addressIDStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid address ID"})
			return
		}

		if err := database.DeleteUserAddress(db.DB, uint(addressID)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting address"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Address deleted successfully"})
	}
}
