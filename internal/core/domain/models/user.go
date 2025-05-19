package models

import "time"

// User represents the structure of a user in the application and in the database
type User struct {
	ID        int
	Name      string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}
