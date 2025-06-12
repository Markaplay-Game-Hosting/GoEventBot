package main

import (
	"context"
	"fmt"
	"github.com/Markaplay-Game-Hosting/GoEventBot/internal/data"
	"github.com/redis/go-redis/v9"
	"sync"

	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
	"log/slog"
	"os"
	"time"
)

type Config struct {
	webhook    string
	recipient  string
	service    *calendar.Service
	calendarId string
	timeSpan   time.Duration
	polling    time.Duration
}

type application struct {
	config Config
	logger *slog.Logger
	models data.Models
	db     *redis.Client
	wg     sync.WaitGroup
}

func Timer(eventInfo data.Event, ch chan<- data.Event) {

	duration := time.Until(eventInfo.StartDate)

	timer := time.NewTimer(duration)
	<-timer.C
	ch <- eventInfo
}

func main() {

	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	logger.Info("Starting checker")
	ctx := context.Background()
	service, err := calendar.NewService(ctx, option.WithCredentialsFile("/config/service_account.json"))
	if err != nil {
		logger.Error("Unable to authenticate to google services")
		panic(err.Error())
	}

	calendarId := os.Getenv("CALENDAR_ID")
	if calendarId == "" {
		logger.Error("environment variable 'CALENDAR_ID' is not set", "error", err.Error())
		panic(err)
	}

	discordUrl := os.Getenv("WEBHOOK_URL")
	if discordUrl == "" {
		logger.Info("environment variable 'WEBHOOK_URL' is not set")
		panic(nil)
	}

	client := ConfigureDB()

	var sendTo = ""
	if os.Getenv("SENDTO") != "" {
		logger.Info(fmt.Sprintf("using %s", os.Getenv("SENDTO")))
		sendTo = os.Getenv("SENDTO")
	}

	timeSpan, err := time.ParseDuration(os.Getenv("TIMESPAN"))
	if err != nil {
		logger.Error("Unable to parse time from environment variable 'TIMESPAN'")
		panic(err.Error())
	}

	pollingRate, err := time.ParseDuration(os.Getenv("POLLING"))
	if err != nil {
		logger.Error("Unable to parse time from environment variable 'POLLING'")
		panic(err.Error())
	}

	app := &application{
		logger: logger,
		models: data.NewModels(client),
		config: Config{
			webhook:    discordUrl,
			recipient:  sendTo,
			service:    service,
			calendarId: calendarId,
			timeSpan:   timeSpan,
			polling:    pollingRate,
		},
		db: client,
	}

	err = app.serve()
	if err != nil {
		logger.Error("Error while serving", err.Error())
		panic(err)
	}

}

func ConfigureDB() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: "",
		DB:       0,
	})
	return client
}
