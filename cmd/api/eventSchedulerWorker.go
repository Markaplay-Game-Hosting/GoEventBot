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
	rrule, err := ParseRRule(event.RRule)
	if err != nil {
		app.logger.Error("Unable to parse RRule", "error", err)
		return nil
	}
	return &eventScheduler{
		JobQueue: make(chan Scheduler),
		Time:     rrule.After(time.Now(), false),
		Event:    event,
	}
}

func (e *eventScheduler) StartScheduleEvent(app *application) {
	app.scheduledEventsTracker = make(map[uuid.UUID]data.Event)
	eventsLoaded <- e.Event.ID
	for {
		delay := time.Until(e.Time)
		if delay <= 0 {
			break
		}
		time.AfterFunc(delay, func() {
			event := <-e.JobQueue
			event.Execute(e.Event)
		})
		// Update e.Time to the next occurrence based on the recurrence rule
		nextRule, err := ParseRRule(e.Event.RRule)
		if err != nil {
			app.logger.Error("Unable to parse RRule", "error", err)
			return
		}
		nextTime := nextRule.After(e.Time, false)
		if nextTime.IsZero() {
			break // Exit loop if no further occurrences
		}
		e.Time = nextTime
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
