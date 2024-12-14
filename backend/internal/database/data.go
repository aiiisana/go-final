package database

import (
	"ecommerce-platform/internal/models"
	"encoding/json"
	"io/ioutil"
	"log"

	"golang.org/x/crypto/bcrypt"
)

func importData[T any](db *Client, filename string, entity []T) {
	fileData, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("Error reading file %s: %v", filename, err)
	}

	err = json.Unmarshal(fileData, &entity)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON data from %s: %v", filename, err)
	}

	for _, item := range entity {
		if err := db.DB.Create(&item).Error; err != nil {
			log.Printf("Error inserting item into database: %v", err)
		}
	}
	log.Printf("%s imported successfully", filename)
}

func importUsers(db *Client, filename string) {
	var users []models.User

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("Error reading file %s: %v", filename, err)
	}

	err = json.Unmarshal(data, &users)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON data: %v", err)
	}

	var userRole models.Role
	if err := db.DB.Where("role_name = ?", "user").First(&userRole).Error; err != nil {
		log.Fatalf("User role not found: %v", err)
	}

	for _, user := range users {
		if len(user.Password) < 6 {
			log.Printf("User %s has invalid password", user.Username)
			continue
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("Error hashing password for user %s: %v", user.Username, err)
			continue
		}

		user.Password = string(hashedPassword)
		user.RoleID = userRole.ID

		if err := db.DB.Create(&user).Error; err != nil {
			log.Printf("Error inserting user %s into database: %v", user.Username, err)
		}
	}

	log.Println("Users imported successfully")
}
