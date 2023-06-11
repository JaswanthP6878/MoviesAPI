package main

import (
	"context"
	"net/http"

	"greenlight.jaswanthp.com/internal/data"
)

type contextkey string

const userContextKey = contextkey("user")

func (app *application) contextSetUser(r *http.Request, user *data.User) *http.Request {
	ctx := context.WithValue(r.Context(), userContextKey, user)
	return r.WithContext(ctx)
}

func (app *application) contextGetUser(r *http.Request) *data.User {
	user, ok := r.Context().Value(userContextKey).(*data.User)
	if !ok {
		panic("missing user details")
	}

	return user
}
