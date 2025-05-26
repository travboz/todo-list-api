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

	return router
}
