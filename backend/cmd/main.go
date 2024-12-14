package main

import (
	"ecommerce-platform/internal/database"
	router "ecommerce-platform/internal/routes"
	"log"
)

func main() {
	database.InitRedis()

	client, err := database.InitializeDatabase()
	if err != nil {
		log.Fatalf("Failed to initialize the database: %v", err)
	}

	r := router.SetupRouter(*client)

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
