package main

import (
	"github.com/Markaplay-Game-Hosting/GoEventBot/internal/data"
	"github.com/Markaplay-Game-Hosting/GoEventBot/internal/validator"
	"net/http"
)

func (app *application) CreateWebhookHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}

	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	webhook := &data.Webhook{
		Name: input.Name,
		URL:  input.URL,
	}

	v := validator.New()
	if data.ValidateWebHook(v, webhook); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	if err := app.models.Webhooks.Insert(webhook); err != nil {
		app.logger.Error("Unable to insert webhook", err.Error())
		app.serverErrorResponse(w, r, err)
		return
	}

	if err := app.writeJSON(w, http.StatusCreated, envelope{"webhook": webhook}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
