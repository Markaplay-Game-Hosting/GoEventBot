package main

import (
	"github.com/Markaplay-Game-Hosting/GoEventBot/internal/data"
	"github.com/google/uuid"
	"time"
)

type eventScheduler struct {
	JobQueue chan Scheduler
	Time     time.Time
	Event    data.Event
}

var eventsLoaded chan uuid.UUID = make(chan uuid.UUID)

func (app *application) NewScheduleEvent(event data.Event) *eventScheduler {

	return &eventScheduler{
		JobQueue: make(chan Scheduler),
		Time:     event.StartDate,
		Event:    event,
	}
}

func (e *eventScheduler) StartScheduleEvent(app *application) {
	app.scheduledEventsTracker = make(map[uuid.UUID]data.Event)
	eventsLoaded <- e.Event.ID
	for {
		delay := time.Until(e.Time)
		time.AfterFunc(delay, func() {
			event := <-e.JobQueue
			event.Execute(e.Event)
		})
	}
}

func (app *application) LoadPendingEvents() {
	app.models.Events.GetAll()
}

func (app *application) GetLoadedEvents() []data.Event {
	var events []data.Event
	for _, event := range app.scheduledEventsTracker {
		events = append(events, event)
	}
	return events
}
