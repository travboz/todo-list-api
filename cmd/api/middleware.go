package main

import (
	"context"
	"crypto/sha256"
	"crypto/subtle"
	"log"
	"net/http"
	"strings"

	"github.com/travboz/backend-projects/todo-list-api/internal/env"
)

func (app *application) basicAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract the username and password from the request
		// Authorization header. If no Authentication header is present
		// or the header value is invalid, then the 'ok' return value
		// will be false.
		username, password, ok := r.BasicAuth()
		if ok {
			// Calculate SHA-256 hashes for the provided and expected
			// usernames and passwords.
			usernameHash := sha256.Sum256([]byte(username))
			passwordHash := sha256.Sum256([]byte(password))

			expectedUsername := env.GetString("AUTH_USERNAME", "travis")
			expectedPassword := env.GetString("AUTH_PASSWORD", "pa55word")

			// why has using sha256? hashed in order to get two equal-length byte slices that can be compared in constant-time.
			expectedUsernameHash := sha256.Sum256([]byte(expectedUsername))
			expectedPasswordHash := sha256.Sum256([]byte(expectedPassword))

			// Use the subtle.ConstantTimeCompare() function to check if
			// the provided username and password hashes equal the
			// expected username and password hashes. ConstantTimeCompare
			// will return 1 if the values are equal, or 0 otherwise.
			// Importantly, we should to do the work to evaluate both the
			// username and password before checking the return values to
			// avoid leaking information.
			usernameMatch := (subtle.ConstantTimeCompare(usernameHash[:], expectedUsernameHash[:]) == 1)
			passwordMatch := (subtle.ConstantTimeCompare(passwordHash[:], expectedPasswordHash[:]) == 1)

			// If username & password are correct, call next handler in chain.
			if usernameMatch && passwordMatch {
				next.ServeHTTP(w, r)
				return
			}
		}

		w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		unauthorisedResponse(app.Logger, w, r)
	})
}

func (app *application) requireToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// read auth header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			bearerUnauthorisedResponse(app.Logger, w, r)
			return
		}

		parts := strings.Split(authHeader, " ")

		if len(parts) != 2 || parts[0] != "Bearer" {
			malformedAuthResponse(app.Logger, w, r, "authorization header is malformed")
			return
		}

		ctx := r.Context()

		token := parts[1]
		valid, err := app.Storage.Tokens.GetAndValidateToken(ctx, token)
		if err != nil || !valid {
			bearerUnauthorisedResponse(app.Logger, w, r)
			return
		}

		// use token to get user id
		user_id, err := app.Storage.Tokens.GetUserIdUsingToken(ctx, token)
		if err != nil {
			serverErrorResponse(app.Logger, w, r, err)
		}

		// Store user ID in context
		ctx = context.WithValue(ctx, contextKeyUserID, user_id)

		// Pass new context down the chain
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				log.Printf("panic: %v\n", err)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
