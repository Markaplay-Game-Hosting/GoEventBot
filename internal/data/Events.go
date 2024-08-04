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

	timeLeft := time.Now().Sub(event.EndDate)
	err := e.DB.Set(ctx, event.ID, nil, timeLeft).Err()
	if err != nil {
		return err
	}
	return nil
}

func (e EventModel) Get(ID string) bool {
	ctx := context.Background()
	eventFound := e.DB.Get(ctx, ID)
	if eventFound == nil {
		return false
	} else {
		return true
	}
}
