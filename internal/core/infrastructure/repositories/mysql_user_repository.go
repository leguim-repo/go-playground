package repository

import (
	"context"
	"database/sql"
	"fmt"
	"go-playground/internal/core/domain/models"
	"go-playground/internal/core/domain/repositories"
	"go-playground/pkg/thelogger"

	_ "github.com/go-sql-driver/mysql" // Driver de MySQL
)

// mysqlUserRepository is an implementation of UserRepository for MySQL
// Note that the db field is a MySQL client passed via dependency injection
// This struct IS UNEXPORTED
type mysqlUserRepository struct {
	db                    *sql.DB
	justAnotherDependency string
	logger                *thelogger.TheLogger
}

// NewMySQLUserRepository creates a new instance of mysqlUserRepository
// This is EXPORTED and can be used by any other package
// In fact this is a constructor
func NewMySQLUserRepository(clientDb *sql.DB, tag string) repositories.UserRepository {
	logger := thelogger.NewTheLogger()

	return &mysqlUserRepository{
		db:                    clientDb,
		justAnotherDependency: tag,
		logger:                logger,
	}
}

// GetByID implements the GetByID method of the UserRepository interface
func (r *mysqlUserRepository) GetByID(ctx context.Context, id int) (*models.User, error) {
	r.logger.Debug("GetByID")
	user := &models.User{}
	query := "SELECT id, name, email, created_at, updated_at FROM users WHERE id = ?"
	// QueryRowContext when we want to obtain only one row
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
func (r *mysqlUserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	r.logger.Debug("GetByEmail")
	user := &models.User{}
	query := "SELECT id, name, email, created_at, updated_at FROM users WHERE email = ?"
	// QueryRowContext when we want to obtain only one row
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
func (r *mysqlUserRepository) Create(ctx context.Context, user *models.User) error {
	query := "INSERT INTO users (name, email, created_at, updated_at) VALUES (?, ?, ?, ?)"
	// ExecContext is used for modify one row for example update insert delete
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
func (r *mysqlUserRepository) Update(ctx context.Context, user *models.User) error {
	query := "UPDATE users SET name = ?, email = ?, updated_at = ? WHERE id = ?"
	// ExecContext is used for modify one row for example update insert delete
	_, err := r.db.ExecContext(ctx, query, user.Name, user.Email, user.UpdatedAt, user.ID)
	if err != nil {
		return fmt.Errorf("error updating user: %w", err)
	}
	return nil
}

// Delete implements the Delete method of the UserRepository interface
func (r *mysqlUserRepository) Delete(ctx context.Context, id int) error {
	r.logger.Debug("Delete")
	query := "DELETE FROM users WHERE id = ?"
	// ExecContext is used for modify one row for example update insert delete
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error deleting user: %w", err)
	}
	return nil
}

func (r *mysqlUserRepository) GetAllUsers(ctx context.Context) ([]*models.User, error) {
	r.logger.Debug("GetAllUsers")
	var users []*models.User

	query := "SELECT * FROM users"
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error querying all users: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		user := &models.User{}
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("error scanning user row: %w", err)
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over user rows: %w", err)
	}

	return users, nil
}
