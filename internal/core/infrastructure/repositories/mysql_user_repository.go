package repository

import (
	"context"
	"database/sql"
	"fmt"
	"go-playground/internal/core/domain/models"

	_ "github.com/go-sql-driver/mysql" // Driver de MySQL
)

// MySQLUserRepository is an implementation of UserRepository for MySQL
type MySQLUserRepository struct {
	db *sql.DB
}

// NewMySQLUserRepository creates a new instance of MySQLUserRepository
func NewMySQLUserRepository(db *sql.DB) *MySQLUserRepository {
	return &MySQLUserRepository{db: db}
}

// GetByID implements the GetByID method of the UserRepository interface
func (r *MySQLUserRepository) GetByID(ctx context.Context, id int) (*models.User, error) {
	user := &models.User{}
	query := "SELECT id, name, email, created_at, updated_at FROM users WHERE id = ?"
	err := r.db.QueryRowContext(ctx, query, id).Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user with ID %d not found", id)
	}
	if err != nil {
		return nil, fmt.Errorf("error getting user by ID: %w", err)
	}
	return user, nil
}

// GetByEmail implements the GetByEmail method of the UserRepository interface
func (r *MySQLUserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	user := &models.User{}
	query := "SELECT id, name, email, created_at, updated_at FROM users WHERE email = ?"
	err := r.db.QueryRowContext(ctx, query, email).Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user with email '%s' not found", email)
	}
	if err != nil {
		return nil, fmt.Errorf("error getting user by email: %w", err)
	}
	return user, nil
}

// Create implements the Create method of the UserRepository interface
func (r *MySQLUserRepository) Create(ctx context.Context, user *models.User) error {
	query := "INSERT INTO users (name, email, created_at, updated_at) VALUES (?, ?, ?, ?)"
	result, err := r.db.ExecContext(ctx, query, user.Name, user.Email, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		return fmt.Errorf("error creating user: %w", err)
	}
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("error getting last insert ID: %w", err)
	}
	user.ID = int(lastInsertID)
	return nil
}

// Update implements the Update method of the UserRepository interface
func (r *MySQLUserRepository) Update(ctx context.Context, user *models.User) error {
	query := "UPDATE users SET name = ?, email = ?, updated_at = ? WHERE id = ?"
	_, err := r.db.ExecContext(ctx, query, user.Name, user.Email, user.UpdatedAt, user.ID)
	if err != nil {
		return fmt.Errorf("error updating user: %w", err)
	}
	return nil
}

// Delete implements the Delete method of the UserRepository interface
func (r *MySQLUserRepository) Delete(ctx context.Context, id int) error {
	query := "DELETE FROM users WHERE id = ?"
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error deleting user: %w", err)
	}
	return nil
}
