package domain

import "context"

// User is main Entity inside user module
type User struct {
	ID        uint
	Username  string
	FirstName string
	LastName  string
}

// UserRepository represents a contract for communication with underlying storage
type UserRepository interface {
	Create(ctx context.Context, user User, password string) (*User, error)
	SearchByUsername(ctx context.Context, username string) ([]User, error)
}
