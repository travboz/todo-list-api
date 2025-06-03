package data

import (
	"time"

	"github.com/travboz/backend-projects/todo-list-api/internal/validator"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Task struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Owner       string             `bson:"owner" json:"owner"` // use a dummy, so just start on users then create tasks and use the first user you insert's id
	Title       string             `bson:"title" json:"title"`
	Description string             `bson:"description" json:"description"`
	Completed   bool               `bson:"completed" json:"completed"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
}

func ValiateTask(v *validator.Validator, task *Task) {
	v.Check(task.Title != "", "title", "must be provided")
	v.Check(len(task.Title) <= 100, "title", "must not be more than 100 bytes long")

	v.Check(task.Description != "", "description", "must be provided")
	v.Check(len(task.Description) <= 1000, "description", "must not be more than 1000 bytes long")
}
