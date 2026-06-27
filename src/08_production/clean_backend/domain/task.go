package domain

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	CollectionTask = "tasks"
)

// Task represents a user's task. Each task belongs to exactly one user.
type Task struct {
	ID     primitive.ObjectID `bson:"_id" json:"-"`
	Title  string             `bson:"title" form:"title" binding:"required" json:"title"`
	UserID primitive.ObjectID `bson:"userID" json:"-"`
}

// TaskRepository defines data access for tasks.
type TaskRepository interface {
	Create(c context.Context, task *Task) error
	FetchByUserID(c context.Context, userID string) ([]Task, error)
}

// TaskUsecase defines business operations for tasks.
type TaskUsecase interface {
	Create(c context.Context, task *Task) error
	FetchByUserID(c context.Context, userID string) ([]Task, error)
}
