package data

import (
	"database/sql"
	"sync"
	"time"
)

type Scheduler struct {
	mu     sync.Mutex
	timers map[int]*time.Timer
	events map[int]Event
	db     *sql.DB
}
