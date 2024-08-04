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
	"math/rand"
	"os"
	"time"
)

type Config struct {
	webhook    string
	recipient  string
	service    *calendar.Service
	calendarId string
}

type application struct {
	config Config
	logger *slog.Logger
	models data.Models
	db     *redis.Client
	wg     sync.WaitGroup
}

func (app *application) GetEvents() (*calendar.Events, error) {
	cal := app.config.service.Events.List(app.config.calendarId)
	rand.Int()
	timeMin := time.Now().Format(time.RFC3339)

	timeMax := time.Now().Add(time.Minute * 5).Format(time.RFC3339)

	cal.TimeMin(timeMin)
	cal.TimeMax(timeMax)

	events, err := cal.Do()
	if err != nil {
		return nil, err
	}
	return events, nil
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
	service, err := calendar.NewService(ctx, option.WithCredentialsFile("service_account.json"))
	if err != nil {
		logger.Error("Unable to authenticate to google services")
		panic(err.Error())
	}

	calendarId := os.Getenv("CALENDAR_ID")
	if calendarId == "" {
		logger.Error("environment variable 'CALENDAR_ID' is not set", err.Error())
		panic(err)
	}

	discordUrl := os.Getenv("WEBHOOK_URL")
	if discordUrl == "" {
		logger.Info("environment variable 'WEBHOOK_URL' is not set")
		panic(nil)
	}

	client := ConfigureDB()

	var sendTo = "@everyone"
	if os.Getenv("SENDTO") != "" {
		logger.Info(fmt.Sprintf("using %s", os.Getenv("SENDTO")))
		sendTo = os.Getenv("SENDTO")
	}

	app := &application{
		logger: logger,
		models: data.NewModels(client),
		config: Config{
			webhook:    discordUrl,
			recipient:  sendTo,
			service:    service,
			calendarId: calendarId,
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