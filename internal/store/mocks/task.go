package mocks

import (
	"context"
	"time"

	"github.com/travboz/backend-projects/todo-list-api/internal/data"
	appErrors "github.com/travboz/backend-projects/todo-list-api/internal/errors"
)

type Task struct {
	Owner       string    `bson:"owner" json:"owner"` // use a dummy, so just start on users then create tasks and use the first user you insert's id
	Title       string    `bson:"title" json:"title"`
	Description string    `bson:"description" json:"description"`
	Completed   bool      `bson:"completed" json:"completed"`
	CreatedAt   time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
}

var mockTask = data.Task{
	Owner:       "trav",
	Title:       "Green Eggs & Mock",
	Description: "Green Mocks & Ham",
	Completed:   false,
	CreatedAt:   time.Now(),
}

type TasksStoreMock struct{}

func (t *TasksStoreMock) Insert(ctx context.Context, task *data.Task) error {
	return nil
}

func (t *TasksStoreMock) GetTaskById(ctx context.Context, id string) (*data.Task, error) {
	switch id {
	case "idone":
		return &mockTask, nil
	default:
		return &data.Task{}, appErrors.ErrRecordNotFound
	}
}

func (t *TasksStoreMock) FetchAllTasks(ctx context.Context, p data.Filters, search string) ([]*data.Task, data.Metadata, error) {
	return []*data.Task{&mockTask}, data.Metadata{CurrentPage: 1, PageSize: 10, TotalRecords: 1}, nil
}

// TODO: sort out testing of update
func (t *TasksStoreMock) UpdateTask(ctx context.Context, id string, task *data.Task) (*data.Task, error) {
	return &mockTask, nil
}

func (t *TasksStoreMock) DeleteTask(ctx context.Context, id string) error {
	return nil
}

func (t *TasksStoreMock) CompleteTask(ctx context.Context, id string) (*data.Task, error) {
	completedMock := &mockTask
	completedMock.Completed = true

	return completedMock, nil
}

func (t *TasksStoreMock) GetTaskOwnerId(ctx context.Context, id string) (string, error) {
	return "trav", nil
}
