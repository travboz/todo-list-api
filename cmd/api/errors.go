package main

import (
	"log/slog"
	"net/http"
)

func logError(logger *slog.Logger, r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
	)

	logger.Error(err.Error(), "method", method, "uri", uri)
}

func errorResponse(logger *slog.Logger, w http.ResponseWriter, r *http.Request, status int, message any) {
	env := envelope{"error": message}

	err := writeJSON(w, status, env, nil)
	if err != nil {
		logError(logger, r, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func serverErrorResponse(logger *slog.Logger, w http.ResponseWriter, r *http.Request, err error) {
	logError(logger, r, err)

	message := "the server encountered a problem and could not process your request"
	errorResponse(logger, w, r, http.StatusInternalServerError, message)
}

func badRequestResponse(logger *slog.Logger, w http.ResponseWriter, r *http.Request, err error) {
	errorResponse(logger, w, r, http.StatusBadRequest, err.Error())
}

func failedValidationResponse(logger *slog.Logger, w http.ResponseWriter, r *http.Request, errors map[string]string) {
	errorResponse(logger, w, r, http.StatusUnprocessableEntity, errors)
}

func notFoundResponse(logger *slog.Logger, w http.ResponseWriter, r *http.Request) {
	message := "the requested resource could not be found"
	errorResponse(logger, w, r, http.StatusNotFound, message)
}

func unauthorisedResponse(logger *slog.Logger, w http.ResponseWriter, r *http.Request) {
	message := "you do not have the correct credentials to access this resource"
	errorResponse(logger, w, r, http.StatusUnauthorized, message)
}

func malformedAuthResponse(logger *slog.Logger, w http.ResponseWriter, r *http.Request, message any) {
	errorResponse(logger, w, r, http.StatusUnauthorized, message)
}

func bearerUnauthorisedResponse(logger *slog.Logger, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("WWW-Authenticate", `Bearer realm="restricted", charset="UTF-8"`)
	unauthorisedResponse(logger, w, r)
}

func forbiddenResponse(logger *slog.Logger, w http.ResponseWriter, r *http.Request) {
	errorResponse(logger, w, r, http.StatusForbidden, "forbidden action")
}
