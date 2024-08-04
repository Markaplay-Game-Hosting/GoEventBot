package data

import "github.com/redis/go-redis/v9"

type Models struct {
	Event EventModel
}

func NewModels(db *redis.Client) Models {
	return Models{
		Event: EventModel{DB: db},
	}
}
