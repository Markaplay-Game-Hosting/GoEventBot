package main

import (
	"encoding/json"
	"errors"
	"net/http"
)

func (app *application) authenticateHandler(w http.ResponseWriter, r *http.Request) {

	state, err := app.models.OAuth.GenerateState()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	url := app.oauth2Config.AuthCodeURL(state)

	http.Redirect(w, r, url, http.StatusFound)
}

func (app *application) callbackHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Optional: validate `state` here if you're storing it

	code := r.URL.Query().Get("code")
	if code == "" {
		app.badRequestResponse(w, r, errors.New("missing authorization code"))
		return
	}

	token, err := app.oauth2Config.Exchange(ctx, code)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	client := app.oauth2Config.Client(ctx, token)

	resp, err := client.Get("https://discord.com/api/users/@me")
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	defer resp.Body.Close()

	var userInfo UserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.logger.Info("Discord user authenticated", "user", userInfo)

	app.setCookieHandler(userInfo, w, r)

	http.Redirect(w, r, "/home", http.StatusFound)
}
