package main

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/travboz/backend-projects/todo-list-api/internal/data"
	"github.com/travboz/backend-projects/todo-list-api/internal/store"
	"github.com/travboz/backend-projects/todo-list-api/internal/validator"
)

func (app *application) createNewTaskHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		owner_id, ok := app.getUserIDFromContext(r.Context())
		if !ok {
			unauthorisedResponse(app.Logger, w, r)
			return
		}

		var input struct {
			Title       string `json:"title"`
			Description string `json:"description"`
		}

		err := readJSON(w, r, &input)
		if err != nil {
			badRequestResponse(app.Logger, w, r, err)
			return
		}

		task := &data.Task{
			Owner:       owner_id,
			Title:       input.Title,
			Description: input.Description,
			Completed:   false,
			CreatedAt:   time.Now(),
		}

		v := validator.New()

		if data.ValidateTask(v, task); !v.Valid() {
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

		err = writeJSON(w, http.StatusCreated, envelope{"task": task}, headers)
		if err != nil {
			serverErrorResponse(app.Logger, w, r, err)
		}
	})
}

func (app *application) fetchAllTasksHandler() http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var input struct {
			data.Pagination
		}

		qs := r.URL.Query()

		v := validator.New()

		input.Pagination.Page = app.readInt(qs, "page", 1, v)
		input.Pagination.PageSize = app.readInt(qs, "page_size", 20, v)

		if data.ValidatePagination(v, input.Pagination); !v.Valid() {
			failedValidationResponse(app.Logger, w, r, v.Errors)
			return
		}

		tasks, err := app.Storage.FetchAllTasks(r.Context(), input.Pagination)
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
		// check that current logged in user is owner of task - as per permission requirements
		ctx := r.Context()
		task_id := readIDParam(r)

		owner_id, err := app.Storage.GetTaskOwnerId(r.Context(), task_id)
		if err != nil {
			switch {
			case errors.Is(err, store.ErrRecordNotFound):
				notFoundResponse(app.Logger, w, r)
			default:
				serverErrorResponse(app.Logger, w, r, err)
			}

			return
		}

		current_user_id, ok := app.getUserIDFromContext(ctx)
		if !ok {
			unauthorisedResponse(app.Logger, w, r)
			return
		}

		if owner_id != current_user_id {
			forbiddenResponse(app.Logger, w, r)
			return
		}

		err = app.Storage.DeleteTask(ctx, task_id)
		if err != nil {
			switch {
			case errors.Is(err, store.ErrRecordNotFound):
				notFoundResponse(app.Logger, w, r)
			default:
				serverErrorResponse(app.Logger, w, r, err)
			}

			return
		}

		err = writeJSON(w, http.StatusOK, envelope{"message": "successful deletion of task with id: " + task_id}, nil)
		if err != nil {
			serverErrorResponse(app.Logger, w, r, err)
		}

	})
}

func (app *application) updateTaskHandler() http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// check that current logged in user is owner of task - as per permission requirements
		ctx := r.Context()
		task_id := readIDParam(r)

		owner_id, err := app.Storage.GetTaskOwnerId(r.Context(), task_id)
		if err != nil {
			switch {
			case errors.Is(err, store.ErrRecordNotFound):
				notFoundResponse(app.Logger, w, r)
			default:
				serverErrorResponse(app.Logger, w, r, err)
			}

			return
		}

		current_user_id, ok := app.getUserIDFromContext(ctx)
		if !ok {
			unauthorisedResponse(app.Logger, w, r)
			return
		}

		if owner_id != current_user_id {
			forbiddenResponse(app.Logger, w, r)
			return
		}

		var input struct {
			Title       string `json:"title"`
			Description string `json:"description"`
			Completed   bool   `json:"completed"`
		}

		err = readJSON(w, r, &input)
		if err != nil {
			badRequestResponse(app.Logger, w, r, err)
			return
		}

		task := &data.Task{
			Title:       input.Title,
			Description: input.Description,
			Completed:   input.Completed,
		}

		v := validator.New()

		if data.ValidateUpdateTask(v, task); !v.Valid() {
			failedValidationResponse(app.Logger, w, r, v.Errors)
			return
		}

		updated, err := app.Storage.UpdateTask(r.Context(), task_id, task)
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

		// check that current logged in user is owner of task - as per permission requirements
		ctx := r.Context()
		task_id := readIDParam(r)

		owner_id, err := app.Storage.GetTaskOwnerId(r.Context(), task_id)
		if err != nil {
			switch {
			case errors.Is(err, store.ErrRecordNotFound):
				notFoundResponse(app.Logger, w, r)
			default:
				serverErrorResponse(app.Logger, w, r, err)
			}

			return
		}

		current_user_id, ok := app.getUserIDFromContext(ctx)
		if !ok {
			unauthorisedResponse(app.Logger, w, r)
			return
		}

		if owner_id != current_user_id {
			forbiddenResponse(app.Logger, w, r)
			return
		}

		completed, err := app.Storage.CompleteTask(r.Context(), task_id)
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
