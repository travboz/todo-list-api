package main

import (
	"log/slog"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func routes(logger *slog.Logger) http.Handler {
	router := httprouter.New()

	router.Handler(http.MethodGet, "/api/v1/healthcheck", healthcheckHandler(logger))

	return router
}
