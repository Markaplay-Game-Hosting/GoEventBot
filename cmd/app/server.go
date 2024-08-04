package main

import (
	"fmt"
	"github.com/Markaplay-Game-Hosting/GoEventBot/internal/data"
	"time"
)

func (app *application) serve() error {
	sendStatusChannel := make(chan data.Event)

	go func() {
		for {
			events, err := app.GetEvents()
			if err != nil {
				app.logger.Error("Unable to get calendar events: ", err.Error())
			}
			if len(events.Items) > 0 {
				for _, event := range events.Items {
					app.logger.Info(fmt.Sprintf("event summary: %s", event.Summary))
					app.logger.Info(fmt.Sprintf("event description: %s", event.Description))
					app.logger.Info(fmt.Sprintf("event start date: %s", event.Start.DateTime))
					eventExist := app.models.Event.Get(event.Id)
					if eventExist == false {
						app.logger.Info("no records found in db, adding it!")

						startDate, err := time.Parse(time.RFC3339, event.Start.DateTime)
						if err != nil {
							app.logger.Error("error while converting start date to time.Time", err.Error())
						}
						endDate, err := time.Parse(time.RFC3339, event.End.DateTime)
						if err != nil {
							app.logger.Error("error while converting end date to time.Time", err.Error())
						}

						eventToCheck := data.Event{
							ID:          event.Id,
							Title:       event.Summary,
							Description: event.Description,
							StartDate:   startDate,
							EndDate:     endDate,
						}
						err = app.models.Event.Insert(&eventToCheck)
						if err != nil {
							app.logger.Error("Unable to add event to DB!")
						} else {
							app.logger.Info("Added event to the DB!")
						}

						Timer(eventToCheck, sendStatusChannel)
					} else {
						app.logger.Info("Event already registered in the DB!")
					}
				}
			}
		}
	}()

	go func() {
		for {
			for eventInfo := range sendStatusChannel {
				msg := FormatMessage(eventInfo)
				err := app.SendMessage(msg)
				if err != nil {
					app.logger.Error("Unable to send message", err.Error())
				} else {
					app.logger.Info("Message sent successfully")
				}
			}
		}
	}()
	return nil
}
