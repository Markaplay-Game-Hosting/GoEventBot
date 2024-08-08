package main

import (
	"fmt"
	"github.com/Markaplay-Game-Hosting/GoEventBot/internal/data"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func (app *application) serve() error {
	sendStatusChannel := make(chan data.Event)

	shutdownError := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit

		app.logger.Info("caught signal", map[string]string{
			"signal": s.String(),
		})

		app.wg.Wait()
		shutdownError <- nil
	}()

	go func() {
		for {
			app.logger.Info(fmt.Sprintf("Waiting %s to check events", app.config.polling))
			time.Sleep(app.config.polling)
			events, err := app.GetEvents()
			if err != nil {
				app.logger.Error("Unable to get calendar events: ", err.Error())
			}
			if len(events.Items) > 0 {
				for _, event := range events.Items {

					startDate, err := time.Parse(time.RFC3339, event.Start.DateTime)
					if err != nil {
						app.logger.Error("error while converting start date to time.Time", err.Error())
					}
					endDate, err := time.Parse(time.RFC3339, event.End.DateTime)
					if err != nil {
						app.logger.Error("error while converting end date to time.Time", err.Error())
					}

					now := time.Now().UTC()

					eventToCheck := data.Event{
						ID:          event.Id,
						Title:       event.Summary,
						Description: event.Description,
						StartDate:   time.Date(now.Year(), now.Month(), now.Day(), startDate.Hour(), startDate.Minute(), 0, 0, time.UTC),
						EndDate:     time.Date(now.Year(), now.Month(), now.Day(), endDate.Hour(), endDate.Minute(), 0, 0, time.UTC),
					}

					eventExist, eventStartDate := app.models.Event.Get(event.Id)

					nowDiff := eventToCheck.StartDate.Sub(now)
					if eventExist == false && nowDiff > 0 && eventStartDate != eventToCheck.StartDate {
						app.logger.Info("no records found in db, adding it!")

						app.logger.Info(fmt.Sprintf("event ID: %s", event.Id))
						app.logger.Info(fmt.Sprintf("event summary: %s", eventToCheck.Title))
						app.logger.Info(fmt.Sprintf("event description: %s", eventToCheck.Description))
						app.logger.Info(fmt.Sprintf("event start date: %s", eventToCheck.StartDate))
						app.logger.Info(fmt.Sprintf("event end date: %s", eventToCheck.EndDate))

						err = app.models.Event.Insert(&eventToCheck)
						if err != nil {
							app.logger.Error("Unable to add event to DB!")
						} else {
							app.logger.Info("Added event to the DB!")
						}

						go Timer(eventToCheck, sendStatusChannel)
					} else {
						if eventExist == true {
							app.logger.Info("Event already registered in the DB")
						} else if nowDiff < 0 {
							app.logger.Info("Event already passed or is ongoing")
						} else if eventStartDate == eventToCheck.StartDate {
							app.logger.Info("Event start date didn't changed")
						}
					}
				}
			}
		}
	}()

	go func() {
		for eventInfo := range sendStatusChannel {
			msg := FormatMessage(eventInfo)
			err := app.SendMessage(msg, eventInfo.Title)
			if err != nil {
				app.logger.Error("Unable to send message", err.Error())
			} else {
				app.logger.Info("Message sent successfully")
			}
		}
	}()

	err := <-shutdownError
	if err != nil {
		return err
	}
	err = app.models.Event.DeleteAll()
	if err != nil {
		return err
	}
	return nil
}
