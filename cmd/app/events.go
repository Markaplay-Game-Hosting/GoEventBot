package main

import (
	"google.golang.org/api/calendar/v3"
	"time"
)

func (app *application) GetEvents() (*calendar.Events, error) {
	cal := app.config.service.Events.List(app.config.calendarId)

	timeMin := time.Now().Format(time.RFC3339)

	timeMax := time.Now().Add(app.config.timeSpan).Format(time.RFC3339)

	cal.TimeMin(timeMin)
	cal.TimeMax(timeMax)

	events, err := cal.Do()
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
