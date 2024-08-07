package main

import (
	"context"
	"google.golang.org/api/calendar/v3"
	"time"
)

func (app *application) GetEvents() (*calendar.Events, error) {

	ctx := context.Background()
	timeMin := time.Now().UTC().Format(time.RFC3339)

	timeMax := time.Now().Add(app.config.timeSpan).UTC().Format(time.RFC3339)

	events, err := app.config.service.Events.List(app.config.calendarId).
		TimeMin(timeMin).
		TimeMax(timeMax).
		ShowDeleted(false).
		Context(ctx).Do()

	if err != nil {
		return nil, err
	}
	return events, nil

}

func (app *application) GetEvent(Id string) (*calendar.Events, error) {
	cal := app.config.service.Events.List(app.config.calendarId)

	timeMin := time.Now().Format(time.RFC3339)

	cal.TimeMin(timeMin)
	cal.ICalUID(Id)

	events, err := cal.Do()
	if err != nil {
		return nil, err
	}
	return events, nil
}
