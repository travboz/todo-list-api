package main

import (
	"errors"
	"net/http"

	"github.com/travboz/backend-projects/todo-list-api/internal/data"
	"github.com/travboz/backend-projects/todo-list-api/internal/store"
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

func (app *application) getUserByIdHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := readIDParam(r)

		current_user_id, ok := app.getUserIDFromContext(r.Context())
		if !ok || current_user_id != id {
			forbiddenResponse(app.Logger, w, r)
			return
		}

		app.Logger.Info("current user id is:", "user_id", current_user_id)

		user, err := app.Storage.UsersModel.Get(id)
		if err != nil {
			switch {
			case errors.Is(err, store.ErrRecordNotFound):
				notFoundResponse(app.Logger, w, r)
			default:
				serverErrorResponse(app.Logger, w, r, err)
			}

			return
		}

		err = writeJSON(w, http.StatusOK, envelope{"user": user}, nil)
		if err != nil {
			serverErrorResponse(app.Logger, w, r, err)
		}
	})
}

func (app *application) userLoginHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var input struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		err := readJSON(w, r, &input)
		if err != nil {
			badRequestResponse(app.Logger, w, r, err)
			return
		}

		user := &data.User{
			Email:    input.Email,
			Password: input.Password,
		}

		id, err := app.Storage.UsersModel.Authenticate(user.Email, user.Password)
		if err != nil {
			serverErrorResponse(app.Logger, w, r, err)
			return
		}

		var token string

		if id != "" {
			// insert a token into the tokens table with an expiry of 1 hour
			token, err = app.Storage.TokensModel.InsertToken(r.Context(), id)
			if err != nil {
				serverErrorResponse(app.Logger, w, r, err)
				return
			}
		}

		err = writeJSON(w, http.StatusCreated, envelope{"login": "successful", "token": token}, nil)
		if err != nil {
			serverErrorResponse(app.Logger, w, r, err)
		}

	})
}
