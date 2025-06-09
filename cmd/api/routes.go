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
	router.Handler(http.MethodGet, "/api/v1/users/:id", app.requireToken(app.getUserByIdHandler()))
	router.Handler(http.MethodPost, "/api/v1/users/login", app.userLoginHandler())

	// task related
	router.Handler(http.MethodPost, "/api/v1/tasks/create", app.requireToken(app.createNewTaskHandler()))
	router.Handler(http.MethodGet, "/api/v1/tasks", app.requireToken(app.fetchAllTasksHandler()))
	router.Handler(http.MethodGet, "/api/v1/tasks/:id", app.requireToken(app.getTasksByIDHandler()))
	router.Handler(http.MethodPut, "/api/v1/tasks/:id/complete", app.requireToken(app.completeTaskHandler()))
	router.Handler(http.MethodPatch, "/api/v1/tasks/:id", app.requireToken(app.updateTaskHandler()))
	router.Handler(http.MethodDelete, "/api/v1/tasks/:id", app.requireToken(app.deleteTaskHandler()))

	return app.recoverPanic(router)
}
