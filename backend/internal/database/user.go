package database

import (
	"ecommerce-platform/internal/models"
	"errors"
	"log"

	"gorm.io/gorm"
)

func GetUserByID(db *gorm.DB, userID uint, currentUser *models.User) (*models.User, error) {
	if currentUser.Role.RoleName != "admin" && currentUser.ID != userID {
		return nil, errors.New("you can only access your own data or be an admin to view others")
	}

	var user models.User
	if err := db.Preload("Role").Where("user_id = ?", userID).First(&user).Error; err != nil {
		log.Println("Error fetching user:", err)
		return nil, err
	}

	if user.ID == 0 {
		return nil, errors.New("user not found")
	}

	return &user, nil
}

func GetAllUsers(db *gorm.DB, currentUser *models.User) ([]models.User, error) {
	if currentUser.Role.ID == 0 || currentUser.Role.RoleName != "admin" {
		log.Printf("Unauthorized access attempt by user ID %d with role %s\n", currentUser.ID, currentUser.Role.RoleName)
		return nil, errors.New("only admins can access all users")
	}

	var users []models.User
	if err := db.Preload("Role").Find(&users).Error; err != nil {
		log.Println("Error fetching all users:", err)
		return nil, err
	}

	return users, nil
}

func CreateUser(db *gorm.DB, user *models.User) error {
	if err := db.Create(&user).Error; err != nil {
		return err
	}
	return nil
}

func UpdateUser(db *gorm.DB, userID uint, updatedData models.User, currentUser *models.User) error {
	if currentUser.Role.RoleName != "admin" && currentUser.ID != userID {
		return errors.New("you can only update your own profile or be an admin to update others")
	}
	if err := db.Model(&models.User{}).Where("id = ?", userID).Updates(updatedData).Error; err != nil {
		log.Println("Error updating user:", err)
		return err
	}
	return nil
}

func DeleteUser(db *gorm.DB, userID uint, currentUser *models.User) error {
	if currentUser.Role.RoleName != "admin" && currentUser.ID != userID {
		return errors.New("you can only delete your own account or be an admin to delete others")
	}
	if err := db.Delete(&models.User{}, "id = ?", userID).Error; err != nil {
		log.Println("Error deleting user:", err)
		return err
	}
	return nil
}

func CreateUserAddress(db *gorm.DB, userID uint, street, city, state, zipCode string) (*models.UserAddress, error) {
	address := models.UserAddress{
		UserID:  userID,
		Street:  street,
		City:    city,
		State:   state,
		ZipCode: zipCode,
	}

	if err := db.Create(&address).Error; err != nil {
		return nil, err
	}
	return &address, nil
}

func GetUserAddresses(db *gorm.DB, userID uint) ([]models.UserAddress, error) {
	var addresses []models.UserAddress
	if err := db.Where("user_id = ?", userID).Find(&addresses).Error; err != nil {
		return nil, err
	}
	return addresses, nil
}

func UpdateUserAddress(db *gorm.DB, addressID uint, updatedData models.UserAddress) error {
	if err := db.Model(&models.UserAddress{}).Where("id = ?", addressID).Updates(updatedData).Error; err != nil {
		return err
	}
	return nil
}

func DeleteUserAddress(db *gorm.DB, addressID uint) error {
	if err := db.Delete(&models.UserAddress{}, "id = ?", addressID).Error; err != nil {
		return err
	}
	return nil
}
