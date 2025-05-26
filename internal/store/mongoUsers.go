package store

import (
	"github.com/travboz/backend-projects/todo-list-api/internal/data"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoStoreUsers struct {
	users *mongo.Collection
}

func (u MongoStoreUsers) Insert(name, email, password string) error {
	return nil
}

func (u MongoStoreUsers) Exists(id string) (bool, error) {
	return false, nil
}

func (u MongoStoreUsers) Get(id string) (*data.User, error) {
	return nil, nil
}
