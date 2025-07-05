package main

import (
	"github.com/Markaplay-Game-Hosting/GoEventBot/cmd/api/models"
	"github.com/Markaplay-Game-Hosting/GoEventBot/internal/data"
	"github.com/Markaplay-Game-Hosting/GoEventBot/internal/validator"
	duration "github.com/channelmeter/iso8601duration"
	"net/http"
	"time"
)

// createEventHandler
// @Summary      Create an event
// @Description  Create a new event
// @Tags         Events
// @Produce      json
// @Param        event 	  body    models.CreateEventRequest   true  "create an event"
// @Success      201
// @Router       /events [post]
// @Security ApiKeyAuth
func (app *application) createEventHandler(w http.ResponseWriter, r *http.Request) {

	var input models.CreateEventRequest

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
		ChannelID:   input.ChannelID,
		GuildID:     input.GuildID,
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

// getEventHandler
// @Summary      Get Event Information
// @Description  Get event information
// @Tags         Events
// @Produce      json
// @param id path string true "event id"
// @Success      200 {object} data.Event
// @Router       /events/{id} [get]
func (app *application) getEventHandler(w http.ResponseWriter, r *http.Request) {
	eventID, err := app.readIDParam(r)
	if err != nil {
		http.Error(w, "Invalid event ID format", http.StatusBadRequest)
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

// getAllEventInstancesHandler
// @Summary      List Upcoming Events
// @Description  List all upcoming events
// @Tags         Events
// @Produce      json
// @Success      200 {array} data.EventInstance
// @Router       /calendar [get]
func (app *application) getAllEventInstancesHandler(w http.ResponseWriter, r *http.Request) {
	events, err := app.models.Events.GetAll()
	if err != nil {
		app.logger.Error("Unable to get all events", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	var eventInstances []data.EventInstance
	for _, event := range events {
		upcoming, err := data.ParseRRule(event.RRule)
		if err != nil {
			app.logger.Error("Unable to parse RRule", err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		now := time.Now()
		firstDayMonth := now.AddDate(0, 0, -now.Day()+1)
		lastDayMonth := now.AddDate(0, 1, -now.Day())
		for _, u := range upcoming.Between(firstDayMonth, lastDayMonth, true) {
			perEventDuration, err := duration.FromString(event.Duration)
			if err != nil {
				app.logger.Error("Unable to parse duration", err.Error())
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			var instance data.EventInstance
			instance.EventID = event.ID
			instance.Title = event.Title
			instance.Description = event.Description
			instance.Duration = event.Duration
			instance.StartDate = u
			instance.EndDate = u.Add(perEventDuration.ToDuration())

			eventInstances = append(eventInstances, instance)
		}
	}

	if err := app.writeJSON(w, http.StatusOK, envelope{"events": eventInstances}, nil); err != nil {
		app.logger.Error("Unable to write JSON", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// getAllEventsHandler
// @Summary      List Events
// @Description  List all upcoming events
// @Tags         Events
// @Produce      json
// @Success      200 {array} data.EventInstance
// @Router       /events [get]
func (app *application) getAllEventsHandler(w http.ResponseWriter, r *http.Request) {
	events, err := app.models.Events.GetAll()
	if err != nil {
		app.logger.Error("Unable to get all events", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if err := app.writeJSON(w, http.StatusOK, envelope{"events": events}, nil); err != nil {
		app.logger.Error("Unable to write JSON", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// deleteEventHandler
// @Summary      Delete Event
// @Description  Delete an event
// @Tags         Events
// @param	id path string true "event id"
// @Produce      json
// @Success      204
// @Router       /events/{id} [delete]
// @Security ApiKeyAuth
func (app *application) deleteEventHandler(w http.ResponseWriter, r *http.Request) {
	eventID, err := app.readIDParam(r)
	if err != nil {
		http.Error(w, "Invalid Event ID format", http.StatusBadRequest)
		return
	}

	err = app.models.Events.Delete(eventID)
	if err != nil {
		app.logger.Error("Unable to delete event", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// updateEventHandler
// @Summary      Update Event
// @Description  Update an event
// @Tags         Events
// @param	id path string true "event id"
// @param event body models.UpdateEventRequest  true  "update an event"
// @Produce      json
// @Success      200
// @Router       /events/{id} [put]
// @Security ApiKeyAuth
func (app *application) updateEventHandler(w http.ResponseWriter, r *http.Request) {
	eventID, err := app.readIDParam(r)
	if err != nil {
		http.Error(w, "Invalid Event ID format", http.StatusBadRequest)
		return
	}

	var input models.UpdateEventRequest

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
	if input.ChannelID != "" {
		event.ChannelID = input.ChannelID
	}
	if input.GuildID != "" {
		event.GuildID = input.GuildID
	}

	if err := app.models.Events.Update(&event); err != nil {
		app.logger.Error("Unable to update event", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
