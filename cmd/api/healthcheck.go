package main

import (
	"log/slog"
	"net/http"
)

func healthcheckHandler(logger *slog.Logger) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data := envelope{
			"status":      "available",
			"system_info": map[string]string{},
		}

		// Encode to json; if there was an error, we log it and send the client
		// a generic error message.
		err := writeJSON(w, http.StatusOK, data, nil)
		if err != nil {
			serverErrorResponse(logger, w, r, err)
		}
	})
}
