package eventWorker

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/Markaplay-Game-Hosting/GoEventBot/internal/bot"
	"github.com/Markaplay-Game-Hosting/GoEventBot/internal/config"
	"github.com/Markaplay-Game-Hosting/GoEventBot/internal/data"
	duration "github.com/channelmeter/iso8601duration"
	"github.com/google/uuid"
)

type EventWorker struct {
	Models    data.Models
	Config    *config.Config
	Bot       *bot.Bot
	Logger    *slog.Logger
	JobTimers sync.Map
}

func (e *EventWorker) Start(ctx context.Context) {

	scheduler := make(chan data.EventInstance, 100)
	var interval time.Duration
	intervalSetting, err := e.Models.Settings.GetByName("load_interval_time")
	if err != nil {
		e.Logger.Warn("no custom interval found in database, using default value (1 minute)")
		interval = time.Minute
	} else {
		interval, err = time.ParseDuration(intervalSetting.Value)
		if err != nil {
			e.Logger.Warn("Error parsing interval time using default value (1 minute)", "error", err)
			interval = time.Minute
		}
	}
	e.Logger.Info("Interval set to", "interval", interval)
	go func() {
		defer close(scheduler)
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				e.Logger.Info("EventWorker shutting down")
				ticker.Stop()
				return
			case <-ticker.C:
				events, err := e.GetUpComingEvents()
				if err != nil {
					e.Logger.Error("Error while fetching upcoming events", "error", err)
				}
				e.Logger.Info("Fetched upcoming events", "count", len(events))
				for _, event := range events {
					scheduler <- event
				}
			}
		}
	}()

	go func() {
		for event := range scheduler {
			select {
			case <-ctx.Done():
				e.Logger.Info("EventWorker shutting down")
				return
			default:
				err := e.Schedule(event)
				if err != nil {
					e.Logger.Error("Error while scheduling event: ", "error", err)
				}
			}

		}
	}()
}

func (e *EventWorker) Schedule(event data.EventInstance) error {
	var err error
	job := &data.Job{
		ID:            uuid.New(),
		EventId:       event.EventID,
		ExecutionDate: event.StartDate,
		Status:        data.Pending,
	}

	_, err = e.Models.Jobs.Insert(job)
	if err != nil {
		return err
	}

	delay := time.Until(event.StartDate)
	if delay <= 0 {
		return nil
	}
	job.Status = data.Running
	err = e.Models.Jobs.Update(job)
	if err != nil {
		e.Logger.Error("Error while updating job status", "error", err.Error())
		return err
	}
	timer := time.AfterFunc(delay, func() {
		e.JobTimers.Delete(job.ID.String())
		message := bot.FormatMessage(event)
		err = e.Bot.SendMessage(event.ChannelID, message)
		if err != nil {
			job.Status = data.Failed
			e.Logger.Error("Error while sending message", "error", err.Error())
		} else {
			e.Logger.Info("Job finished", "jobID", job.ID)
			job.Status = data.Completed
		}

		if err = e.Models.Jobs.Update(job); err != nil {
			e.Logger.Error("Error while updating job status", "error", err.Error())
		}
	})

	e.JobTimers.Store(job.ID.String(), timer)

	return nil
}

func (e *EventWorker) GetUpComingEvents() ([]data.EventInstance, error) {

	var upcomingEvents []data.EventInstance
	var err error
	events, err := e.Models.Events.GetAll()
	if err != nil {
		return upcomingEvents, err
	}

	for _, event := range events {
		upcoming, err := data.ParseRRule(event.RRule)
		if err != nil {
			return upcomingEvents, err
		}
		now := time.Now()
		for _, u := range upcoming.Between(now, now.Add(time.Minute), true) {
			perEventDuration, err := duration.FromString(event.Duration)
			if err != nil {
				return upcomingEvents, err
			}
			instance := data.EventInstance{
				EventID:     event.ID,
				Title:       event.Title,
				Description: event.Description,
				Duration:    event.Duration,
				ChannelID:   event.ChannelID,
				StartDate:   u,
				EndDate:     u.Add(perEventDuration.ToDuration()),
			}

			upcomingEvents = append(upcomingEvents, instance)
		}
	}

	return upcomingEvents, nil
}

func (e *EventWorker) CancelJob(ID uuid.UUID) {
	if val, ok := e.JobTimers.Load(ID.String()); ok {
		if timer, ok := val.(*time.Timer); ok {
			timer.Stop()
			e.JobTimers.Delete(ID.String())
		}
	}
}
