package store

import "github.com/travboz/backend-projects/todo-list-api/internal/data"

type UsersModel interface {
	Insert(name, email, password string) error
	Exists(id string) (bool, error)
	Get(id string) (*data.User, error)
}

type TasksModel interface{}

type Store interface {
	UsersModel
	TasksModel
}
