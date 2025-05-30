package main

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/travboz/backend-projects/todo-list-api/internal/data"
	"github.com/travboz/backend-projects/todo-list-api/internal/store"
	"github.com/travboz/backend-projects/todo-list-api/internal/validator"
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
			// Owner:       input.Owner,
			Title:       input.Title,
			Description: input.Description,
			Completed:   false,
			CreatedAt:   time.Now(),
		}

		v := validator.New()

		if data.ValiateTask(v, task); !v.Valid() {
			failedValidationResponse(app.Logger, w, r, v.Errors)
			return
		}

		err = app.Storage.TasksModel.Insert(r.Context(), task)
		if err != nil {
			serverErrorResponse(app.Logger, w, r, err)
			return
		}

		// where to find the new article:
		headers := make(http.Header)
		headers.Set("Location", fmt.Sprintf("/api/v1/tasks/%s", task.ID.Hex()))

		err = writeJSON(w, http.StatusCreated, envelope{"task": task}, nil)
		if err != nil {
			serverErrorResponse(app.Logger, w, r, err)
		}
	})
}

// =======================================

func (app *application) fetchAllTasksHandler() http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tasks, err := app.Storage.FetchAllTasks(r.Context())
		if err != nil {
			serverErrorResponse(app.Logger, w, r, err)
			return
		}

		err = writeJSON(w, http.StatusOK, envelope{"tasks": tasks}, nil)
		if err != nil {
			serverErrorResponse(app.Logger, w, r, err)
		}
	})
}

func (app *application) getTasksByIDHandler() http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := readIDParam(r)

		task, err := app.Storage.GetTaskById(r.Context(), id)
		if err != nil {
			switch {
			case errors.Is(err, store.ErrRecordNotFound):
				notFoundResponse(app.Logger, w, r)
			default:
				serverErrorResponse(app.Logger, w, r, err)
			}

			return
		}

		err = writeJSON(w, http.StatusOK, envelope{"task": task}, nil)
		if err != nil {
			serverErrorResponse(app.Logger, w, r, err)
		}

	})
}

func (app *application) deleteTaskHandler() http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := readIDParam(r)

		err := app.Storage.DeleteTask(r.Context(), id)
		if err != nil {
			switch {
			case errors.Is(err, store.ErrRecordNotFound):
				notFoundResponse(app.Logger, w, r)
			default:
				serverErrorResponse(app.Logger, w, r, err)
			}

			return
		}

		err = writeJSON(w, http.StatusOK, envelope{"message": "succesful deletion of task with id: " + id}, nil)
		if err != nil {
			serverErrorResponse(app.Logger, w, r, err)
		}

	})
}

func (app *application) updateTaskHandler() http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := readIDParam(r)

		var input struct {
			Title       string `json:"title" bson:"title"`
			Description string `json:"description" bson:"description"`
			Completed   bool   `json:"completed" bson:"completed"`
		}

		err := readJSON(w, r, &input)
		if err != nil {
			badRequestResponse(app.Logger, w, r, err)
			return
		}

		task := &data.Task{
			Title:       input.Title,
			Description: input.Description,
			Completed:   input.Completed,
		}

		// TODO: Add validation
		// v := validator.New()

		// if data.ValidateArticle(v, article); !v.Valid() {
		// 	failedValidationResponse(logger, w, r, v.Errors)
		// 	return
		// }

		updated, err := app.Storage.UpdateTask(r.Context(), id, task)
		if err != nil {
			switch {
			case errors.Is(err, store.ErrRecordNotFound):
				notFoundResponse(app.Logger, w, r)
			default:
				serverErrorResponse(app.Logger, w, r, err)
			}

			return
		}

		err = writeJSON(w, http.StatusOK, envelope{"task": updated}, nil)
		if err != nil {
			serverErrorResponse(app.Logger, w, r, err)
		}

	})
}

func (app *application) completeTaskHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := readIDParam(r)

		completed, err := app.Storage.CompleteTask(r.Context(), id)
		if err != nil {
			switch {
			case errors.Is(err, store.ErrRecordNotFound):
				notFoundResponse(app.Logger, w, r)
			default:
				serverErrorResponse(app.Logger, w, r, err)
			}

			return
		}

		err = writeJSON(w, http.StatusOK, envelope{"task": completed}, nil)
		if err != nil {
			serverErrorResponse(app.Logger, w, r, err)
		}

	})
}
