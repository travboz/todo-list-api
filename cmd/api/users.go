package main

import (
	"net/http"

	"github.com/travboz/backend-projects/todo-list-api/internal/data"
)

func (app *application) registerNewUserHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var input struct {
			Name     string `json:"name"`
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		err := readJSON(w, r, &input)
		if err != nil {
			badRequestResponse(app.Logger, w, r, err)
			return
		}

		user := &data.User{
			Name:     input.Name,
			Email:    input.Email,
			Password: input.Password,
		}

		err = app.Storage.UsersModel.Insert(user)
		if err != nil {
			serverErrorResponse(app.Logger, w, r, err)
			return
		}

		err = writeJSON(w, http.StatusCreated, envelope{"user successfully created": user}, nil)
		if err != nil {
			serverErrorResponse(app.Logger, w, r, err)
		}
	})
}
