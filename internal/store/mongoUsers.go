package store

import (
	"github.com/travboz/backend-projects/todo-list-api/internal/data"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoDBStoreUsers struct {
	users *mongo.Collection
}

func (u MongoDBStoreUsers) Insert(name, email, password string) error {
	return nil
}

func (u MongoDBStoreUsers) Exists(id string) (bool, error) {
	return false, nil
}

func (u MongoDBStoreUsers) Get(id string) (*data.User, error) {
	return nil, nil
}
