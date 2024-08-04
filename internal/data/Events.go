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

func (e EventModel) Get(ID string) (bool, time.Time, error) {
	ctx := context.Background()
	eventFound, err := e.DB.Get(ctx, ID).Result()
	var date time.Time
	if err != nil {
		return false, date, err
	}
	if eventFound == "" {
		return false, date, nil
	} else {
		date, err := time.Parse(time.RFC3339, eventFound)
		if err != nil {
			return false, date, err
		}
		return true, date, nil
	}
}

func (e EventModel) GetAll() ([]string, error) {
	var events []string
	ctx := context.Background()
	events, err := e.DB.Keys(ctx, "*").Result()
	if err != nil {
		return events, err
	}

	return events, nil
}

func (e EventModel) Delete(ID string) error {
	ctx := context.Background()
	err := e.DB.Del(ctx, ID).Err()
	if err != nil {
		return err
	}
	return nil
}

func (e EventModel) DeleteAll() error {
	ctx := context.Background()
	all, err := e.GetAll()
	if err != nil {
		return err
	}
	for _, key := range all {
		err := e.DB.Del(ctx, key).Err()
		if err != nil {
			return err
		}
	}

	return nil
}
