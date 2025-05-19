package usescases

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"go-playground/internal/core/domain/models"
	"go-playground/internal/core/infrastructure/repositories"
)

func UsecaseMysqlUserRepository() {
	dsn := "root:toor@tcp(127.0.0.1:3306)/database_name?parseTime=true"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Check connection
	err = db.Ping()
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}
	fmt.Println("Successfully connected to MySQL!")

	// Initialize repository
	userRepo := repository.NewMySQLUserRepository(db)

	// Use a base context for the example
	ctx := context.Background()

	// Create a new user
	newUser := &models.User{
		Name:      "John Smith",
		Email:     "john.smith@example.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err = userRepo.Create(ctx, newUser)
	if err != nil {
		log.Fatalf("Error creating user: %v", err)
	}
	fmt.Printf("User created with ID: %d\n", newUser.ID)

	// Get user by ID
	foundUser, err := userRepo.GetByID(ctx, newUser.ID)
	if err != nil {
		log.Fatalf("Error getting user by ID: %v", err)
	}
	fmt.Printf("Found user: %+v\n", foundUser)

	// Update user
	foundUser.Name = "John Smith"
	foundUser.UpdatedAt = time.Now()
	err = userRepo.Update(ctx, foundUser)
	if err != nil {
		log.Fatalf("Error updating user: %v", err)
	}
	fmt.Printf("User updated: %+v\n", foundUser)

	// Get user by email
	foundUserByEmail, err := userRepo.GetByEmail(ctx, "john.smith@example.com")
	if err != nil {
		log.Fatalf("Error getting user by email: %v", err)
	}
	fmt.Printf("Found user by email: %+v\n", foundUserByEmail)

	// Delete user
	err = userRepo.Delete(ctx, newUser.ID)
	if err != nil {
		log.Fatalf("Error deleting user: %v", err)
	}
	fmt.Printf("User with ID %d deleted.\n", newUser.ID)
}
