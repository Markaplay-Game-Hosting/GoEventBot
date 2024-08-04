package data

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

type Event struct {
	ID          string
	Title       string
	Description string
	StartDate   time.Time
	EndDate     time.Time
}

type EventModel struct {
	DB *redis.Client
}

func (e EventModel) Insert(event *Event) error {
	ctx := context.Background()

	timeLeft := event.EndDate.Sub(time.Now())
	err := e.DB.Set(ctx, event.ID, event.StartDate.Format(time.RFC3339), timeLeft).Err()
	if err != nil {
		return err
	}
	return nil
}

func (e EventModel) Get(ID string) (time.Time, error) {
	ctx := context.Background()
	eventFound, err := e.DB.Get(ctx, ID).Result()
	var date time.Time
	if err != nil {
		return date, err
	}
	if eventFound != ID {
		return date, nil
	} else {
		date, err := time.Parse(time.RFC3339, eventFound)
		if err != nil {
			return date, err
		}
		return date, nil
	}
}
