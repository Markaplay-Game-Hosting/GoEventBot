package main

import (
	"github.com/Markaplay-Game-Hosting/GoEventBot/internal/data"
	"github.com/google/uuid"
)

type Scheduler interface {
	Execute(event data.Event)
}

func (app *application) Execute(event data.Event) {
	msg := FormatMessage(event)
	err := app.SendMessage(msg, event.Title, uuid.Nil)
	if err != nil {
		app.logger.Error("Unable to send message", err.Error())
	} else {
		app.logger.Info("Message sent successfully")
	}
}
