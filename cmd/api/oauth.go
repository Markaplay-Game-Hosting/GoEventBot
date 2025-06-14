package main

import (
	"golang.org/x/oauth2"
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

	token, err := app.oauth2Config.Exchange(ctx, r.URL.Query().Get("code"))
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	userInfo, err := app.provider.UserInfo(ctx, oauth2.StaticTokenSource(token))
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	app.logger.Info("User Info", "user", userInfo)

}
