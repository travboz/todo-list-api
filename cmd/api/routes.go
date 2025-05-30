package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.Handler(http.MethodGet, "/api/v1/healthcheck", app.healthcheckHandler())

	// user related
	router.Handler(http.MethodPost, "/api/v1/users/register", app.registerNewUserHandler())

	// task related
	router.Handler(http.MethodPost, "/api/v1/tasks/create", app.createNewTaskHandler())
	router.Handler(http.MethodGet, "/api/v1/tasks", app.fetchAllTasksHandler())
	router.Handler(http.MethodGet, "/api/v1/tasks/:id", app.getTasksByIDHandler())
	router.Handler(http.MethodPut, "/api/v1/tasks/:id/complete", app.completeTaskHandler())
	router.Handler(http.MethodPatch, "/api/v1/tasks/:id", app.updateTaskHandler())
	router.Handler(http.MethodDelete, "/api/v1/tasks/:id", app.deleteTaskHandler())

	return router
}
