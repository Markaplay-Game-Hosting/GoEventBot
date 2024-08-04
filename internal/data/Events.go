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
	err := e.DB.Set(ctx, event.ID, event.StartDate, timeLeft).Err()
	if err != nil {
		return err
	}
	return nil
}

func (e EventModel) Get(ID string) (bool, error) {
	ctx := context.Background()
	eventFound, err := e.DB.Get(ctx, ID).Result()
	if err != nil {
		return false, err
	}
	if eventFound != ID {
		return false, nil
	} else {
		return true, nil
	}
}
