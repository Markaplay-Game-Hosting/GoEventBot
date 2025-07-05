package main

import (
	"context"
	"net/http"

	"github.com/Markaplay-Game-Hosting/GoEventBot/internal/data"
)

type contextKey string

const userContextKey = contextKey("user")

func (app *application) contextSetUser(r *http.Request, user *data.User, oauthUser *UserInfo) *http.Request {
	if user == nil {
		ctx := context.WithValue(r.Context(), userContextKey, oauthUser)
		return r.WithContext(ctx)
	}
	ctx := context.WithValue(r.Context(), userContextKey, user)
	return r.WithContext(ctx)
}

func (app *application) contextGetUser(r *http.Request) (*data.User, *UserInfo) {
	user, uok := r.Context().Value(userContextKey).(*data.User)
	oauthUser, dok := r.Context().Value(userContextKey).(*UserInfo)

	if !uok && !dok {
		panic("user not found")
	}
	return user, oauthUser
}
