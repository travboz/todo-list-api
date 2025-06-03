package main

import "context"

type contextKey string

const contextKeyUserID = contextKey("userID")

func (app *application) getUserIDFromContext(ctx context.Context) (string, bool) {
	user_id, ok := ctx.Value(contextKeyUserID).(string)
	return user_id, ok
}
