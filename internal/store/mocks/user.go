package mocks

import (
	"github.com/travboz/backend-projects/todo-list-api/internal/data"
	"github.com/travboz/backend-projects/todo-list-api/internal/errors"
)

type UsersStoreMock struct {
}

func (u UsersStoreMock) Insert(user *data.User) error {
	switch user.Email {
	case "dupe@example.com":
		return errors.ErrDuplicateEmail
	default:
		return nil
	}
}

func (u UsersStoreMock) Authenticate(email, password string) (string, error) {
	if email == "test@example.com" && password == "pa55word" {
		return "user", nil
	}

	return "", errors.ErrInvalidCredentials
}

func (u UsersStoreMock) Get(id string) (*data.User, error) {
	switch id {
	case "user1":
		return &data.User{
			Name:     "user",
			Email:    "test@example.com",
			Password: "pa55word",
		}, nil
	default:
		return &data.User{}, nil
	}
}
