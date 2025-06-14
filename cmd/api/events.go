package main

import (
	"github.com/Markaplay-Game-Hosting/GoEventBot/internal/data"
	"github.com/Markaplay-Game-Hosting/GoEventBot/internal/validator"
	"github.com/google/uuid"
	"net/http"
)

func (app *application) createEventHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		Title       string    `json:"title"`
		Description string    `json:"description"`
		Duration    string    `json:"duration"`
		RRule       string    `json:"rrule"`
		WebhookId   uuid.UUID `json:"webhook_id"`
	}

	if err := app.readJSON(w, r, &input); err != nil {
		app.logger.Error("Unable to read JSON", err.Error())
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	event := &data.Event{
		Title:       input.Title,
		Description: input.Description,
		Duration:    input.Duration,
		RRule:       input.RRule,
		IsActive:    true,
		WebhookID:   input.WebhookId,
	}

	v := validator.New()
	if data.ValidateEvent(v, event); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	if err := app.models.Events.Insert(event); err != nil {
		app.logger.Error("Unable to insert event", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (app *application) getEventHandler(w http.ResponseWriter, r *http.Request) {
	eventID := r.URL.Query().Get("id")
	if eventID == "" {
		http.Error(w, "Event ID is required", http.StatusBadRequest)
		return
	}

	event, err := app.models.Events.Get(eventID)
	if err != nil {
		app.logger.Error("Unable to get event", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if err := app.writeJSON(w, http.StatusOK, envelope{"event": event}, nil); err != nil {
		app.logger.Error("Unable to write JSON", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (app *application) getAllEventsHandler(w http.ResponseWriter, r *http.Request) {
	events, err := app.models.Events.GetAll()
	if err != nil {
		app.logger.Error("Unable to get all events", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if err := app.writeJSON(w, http.StatusOK, envelope{"events": events}, nil); err != nil {
		app.logger.Error("Unable to write JSON", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (app *application) deleteEventHandler(w http.ResponseWriter, r *http.Request) {
	eventID := r.URL.Query().Get("id")
	if eventID == "" {
		http.Error(w, "Event ID is required", http.StatusBadRequest)
		return
	}

	err := app.models.Events.Delete(eventID)
	if err != nil {
		app.logger.Error("Unable to delete event", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (app *application) updateEventHandler(w http.ResponseWriter, r *http.Request) {
	eventID := r.URL.Query().Get("id")
	if eventID == "" {
		http.Error(w, "Event ID is required", http.StatusBadRequest)
		return
	}

	var input struct {
		Title       string    `json:"title,omitempty"`
		Description string    `json:"description,omitempty"`
		Duration    string    `json:"duration,omitempty"`
		RRule       string    `json:"rrule,omitempty"`
		IsActive    bool      `json:"is_active,omitempty"`
		WebhookId   uuid.UUID `json:"webhook_id,omitempty"`
	}

	if err := app.readJSON(w, r, &input); err != nil {
		app.logger.Error("Unable to read JSON", err.Error())
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	event, err := app.models.Events.Get(eventID)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	if input.Title != "" {
		event.Title = input.Title
	}
	if input.Description != "" {
		event.Description = input.Description
	}
	if input.Duration != "" {
		event.Duration = input.Duration
	}
	if input.RRule != "" {
		event.RRule = input.RRule
	}
	if input.IsActive != event.IsActive {
		event.IsActive = input.IsActive
	}
	if input.WebhookId != uuid.Nil {
		event.WebhookID = input.WebhookId
	}

	if err := app.models.Events.Update(&event); err != nil {
		app.logger.Error("Unable to update event", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (app *application) getActiveEventsHandler(w http.ResponseWriter, r *http.Request) {
	events, err := app.models.Events.GetActiveEvents()
	if err != nil {
		app.logger.Error("Unable to get active events", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if err := app.writeJSON(w, http.StatusOK, envelope{"active_events": events}, nil); err != nil {
		app.logger.Error("Unable to write JSON", err.Error())
		app.serverErrorResponse(w, r, err)
	}
}
