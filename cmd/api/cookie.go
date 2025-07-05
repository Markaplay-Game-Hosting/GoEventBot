package main

import (
	"errors"
	"fmt"
	"net/http"
)

const (
	cookieName = "GoEvent"
)

type UserInfo struct {
	ID            string `json:"id"`
	Username      string `json:"username"`
	Discriminator string `json:"discriminator"`
	Email         string `json:"email"`
}

func (app *application) setCookieHandler(user UserInfo, w http.ResponseWriter, r *http.Request) {
	encoded, err := app.cookieHandler.Encode(cookieName, user)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	cookie := &http.Cookie{
		Name:     cookieName,
		Value:    encoded,
		Path:     "/",
		MaxAge:   1800,
		Secure:   false,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, cookie)

}

func (app *application) getCookieHandler(r *http.Request) (UserInfo, error) {
	var userInfo UserInfo

	c, err := r.Cookie(cookieName)
	if err != nil {
		// no cookie → caller can treat as "anonymous"
		if errors.Is(err, http.ErrNoCookie) {
			return userInfo, http.ErrNoCookie
		}
		// malformed header or other error
		return userInfo, err
	}

	if err := app.cookieHandler.Decode(cookieName, c.Value, &userInfo); err != nil {
		// invalid signature, tampered payload, etc.
		return userInfo, fmt.Errorf("invalid cookie: %w", err)
	}

	return userInfo, nil
}
