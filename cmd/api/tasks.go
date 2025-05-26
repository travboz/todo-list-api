package main

import (
	"net/http"
	"time"

	"github.com/travboz/backend-projects/todo-list-api/internal/data"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (app *application) createNewTaskHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var input struct {
			Owner       primitive.ObjectID `json:"owner,omitempty"` // use a dummy, so just start on users then create tasks and use the first user you insert's id
			Title       string             `json:"title"`
			Description string             `json:"description"`
		}

		err := readJSON(w, r, &input)
		if err != nil {
			badRequestResponse(app.Logger, w, r, err)
			return
		}

		task := &data.Task{
			Owner:       input.Owner,
			Title:       input.Title,
			Description: input.Description,
			Completed:   false,
			CreatedAt:   time.Now(),
		}

		err = app.Storage.TasksModel.Insert(r.Context(), task)
		if err != nil {
			serverErrorResponse(app.Logger, w, r, err)
			return
		}

		err = writeJSON(w, http.StatusCreated, envelope{"task": task}, nil)
		if err != nil {
			serverErrorResponse(app.Logger, w, r, err)
		}
	})
}
