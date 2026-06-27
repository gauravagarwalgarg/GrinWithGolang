// PostgreSQL Repository Pattern with sqlx, migrations, and connection pooling.
//
// LEARNING NOTES:
// - sqlx extends database/sql with struct scanning and named parameters
// - Migrations: versioned schema changes (up/down) for reproducible deployments
// - Connection pooling: db.SetMaxOpenConns controls concurrent connections
// - Transactions: wrap multiple operations in BEGIN/COMMIT for atomicity
// - Prepared statements: Postgres pre-compiles queries for faster repeated execution
//
// For C++ devs: Think of sqlx as a type-safe ORM-lite (no magic, just struct mapping).
// The connection pool is like a thread-safe object pool with RAII-style cleanup.
//
// Run: go run main.go
// Requires: PostgreSQL on localhost:5432
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// --- Domain Models ---

type User struct {
	ID        int       `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	Email     string    `db:"email" json:"email"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type Order struct {
	ID        int       `db:"id" json:"id"`
	UserID    int       `db:"user_id" json:"user_id"`
	Amount    float64   `db:"amount" json:"amount"`
	Status    string    `db:"status" json:"status"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

// --- Repository Interface (clean architecture) ---

type UserRepository interface {
	Create(ctx context.Context, name, email string) (*User, error)
	GetByID(ctx context.Context, id int) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	List(ctx context.Context, limit, offset int) ([]User, error)
}

// --- Repository Implementation ---

type postgresUserRepo struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) UserRepository {
	return &postgresUserRepo{db: db}
}

func (r *postgresUserRepo) Create(ctx context.Context, name, email string) (*User, error) {
	var user User
	query := `INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id, name, email, created_at`
	err := r.db.QueryRowxContext(ctx, query, name, email).StructScan(&user)
	return &user, err
}

func (r *postgresUserRepo) GetByID(ctx context.Context, id int) (*User, error) {
	var user User
	err := r.db.GetContext(ctx, &user, `SELECT * FROM users WHERE id = $1`, id)
	return &user, err
}

func (r *postgresUserRepo) GetByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	err := r.db.GetContext(ctx, &user, `SELECT * FROM users WHERE email = $1`, email)
	return &user, err
}

func (r *postgresUserRepo) List(ctx context.Context, limit, offset int) ([]User, error) {
	var users []User
	err := r.db.SelectContext(ctx, &users, `SELECT * FROM users ORDER BY id LIMIT $1 OFFSET $2`, limit, offset)
	return users, err
}

// --- Transaction Example ---

func CreateOrderWithTransaction(db *sqlx.DB, ctx context.Context, userID int, amount float64) (*Order, error) {
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	// Defer rollback no-op if commit succeeds
	defer tx.Rollback()

	// Check user exists
	var exists bool
	err = tx.GetContext(ctx, &exists, `SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)`, userID)
	if err != nil || !exists {
		return nil, fmt.Errorf("user %d not found", userID)
	}

	// Insert order
	var order Order
	err = tx.QueryRowxContext(ctx,
		`INSERT INTO orders (user_id, amount, status) VALUES ($1, $2, 'pending') RETURNING *`,
		userID, amount,
	).StructScan(&order)
	if err != nil {
		return nil, err
	}

	// Commit the transaction
	return &order, tx.Commit()
}

// --- Schema Migration (run once) ---

const schema = `
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS orders (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    amount DECIMAL(10,2) NOT NULL,
    status VARCHAR(50) DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_orders_user_id ON orders(user_id);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
`

func main() {
	// Connection string (use env vars in production)
	dsn := "host=localhost port=5432 user=postgres password=postgres dbname=learning_db sslmode=disable"

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer db.Close()

	// Connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Run migrations
	db.MustExec(schema)
	fmt.Println("Schema migrated successfully")

	// --- Demo ---
	ctx := context.Background()
	repo := NewUserRepository(db)

	// Create a user
	user, err := repo.Create(ctx, "Gaurav", "gaurav@example.com")
	if err != nil {
		log.Printf("Create user (may already exist): %v", err)
		user, _ = repo.GetByEmail(ctx, "gaurav@example.com")
	}
	fmt.Printf("User: %+v\n", user)

	// Create an order with transaction
	order, err := CreateOrderWithTransaction(db, ctx, user.ID, 299.99)
	if err != nil {
		log.Printf("Create order: %v", err)
	} else {
		fmt.Printf("Order: %+v\n", order)
	}

	// List users
	users, _ := repo.List(ctx, 10, 0)
	fmt.Printf("Total users: %d\n", len(users))
}
