// Package domain contains all business entities and interface contracts.
//
// LEARNING NOTES:
// - The domain package has ZERO external dependencies (no gin, no mongo driver imports)
// - It defines WHAT the system does, not HOW (interfaces only)
// - For C++ devs: Think of this as your pure abstract base classes
package domain

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	CollectionUser = "users"
)

// User is the core entity. All auth flows revolve around this model.
type User struct {
	ID       primitive.ObjectID `bson:"_id" json:"id"`
	Name     string             `bson:"name" json:"name"`
	Email    string             `bson:"email" json:"email"`
	Password string             `bson:"password" json:"-"` // Never expose password in JSON
}

// UserRepository defines the data access contract for users.
// The repository layer implements this; the usecase layer depends on it.
type UserRepository interface {
	Create(c context.Context, user *User) error
	Fetch(c context.Context) ([]User, error)
	GetByEmail(c context.Context, email string) (User, error)
	GetByID(c context.Context, id string) (User, error)
}
