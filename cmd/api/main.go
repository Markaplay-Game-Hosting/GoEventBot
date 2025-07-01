package main

import (
	"context"
	"database/sql"
	"errors"
	"expvar"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/Markaplay-Game-Hosting/GoEventBot/internal/config"
	"github.com/Markaplay-Game-Hosting/GoEventBot/internal/data"
	"github.com/Markaplay-Game-Hosting/GoEventBot/internal/vcs"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/golang-migrate/migrate/v4/source/github"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
)

var (
	version = vcs.Version()
)

type application struct {
	config                 config.Config
	logger                 *slog.Logger
	models                 data.Models
	scheduledEventsTracker map[uuid.UUID]data.Event
	oauth2Config           oauth2.Config
	provider               *oidc.Provider
	wg                     sync.WaitGroup
}

// @title           Go Event Bot API
// @version         1.0
// @description     Go Event Bot API for managing events and scheduling.
// @BasePath  /v1/
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	cfg, err := config.Load()

	if err != nil {
		log.Panic("Error loading configuration: ", err)
	}

	displayVersion := flag.Bool("version", false, "Display version and exit")

	flag.Parse()

	if *displayVersion {
		fmt.Printf("Version:\t%s\n", version)
		os.Exit(0)
	}

	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))

	db, err := openDB(*cfg)
	if err != nil {
		logger.Error("Error while opening database connection: ", err)
		os.Exit(1)
	}
	defer func(db *sql.DB) {
		err = db.Close()
	}(db)

	if err != nil {
		logger.Error("Error closing database connection: ", err)
		os.Exit(1)
	}

	logger.Info("database connection pool established")

	logger.Info("running database migrations")
	err = migrateDB(db, logger)
	if err != nil {
		logger.Error("Error running database migrations: ", err)
		os.Exit(1)
	}
	logger.Info("database migrations done")

	oauth2Config, provider, err := setupOauth(*cfg)
	if err != nil {
		logger.Error("Error setting up OAuth2 configuration: ", err)
		os.Exit(1)
	}

	expvar.NewString("version").Set(version)

	expvar.Publish("goroutines", expvar.Func(func() any {
		return runtime.NumGoroutine()
	}))

	expvar.Publish("database", expvar.Func(func() any {
		return db.Stats()
	}))

	expvar.Publish("timestamp", expvar.Func(func() any {
		return time.Now().Unix()
	}))
	app := &application{
		config:                 *cfg,
		logger:                 logger,
		models:                 data.NewModels(db),
		scheduledEventsTracker: make(map[uuid.UUID]data.Event),
		oauth2Config:           oauth2Config,
		provider:               provider,
	}
	err = app.serve()
	if err != nil {
		logger.Error("Error when running the app: ", "error", err)
	}
}

func openDB(cfg config.Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.DB.DSN)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.DB.MaxOpenConns)
	db.SetMaxIdleConns(cfg.DB.MaxIdleConns)

	duration, err := time.ParseDuration(cfg.DB.MaxIdleTime)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func migrateDB(db *sql.DB, logger *slog.Logger) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	logger.Info("Setting up database instance for migration")
	if err != nil {
		return fmt.Errorf("error setting up instance for migration: %w", err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"goeventbotdb", driver)
	if err != nil {
		return fmt.Errorf("error creating new migration instance: %w", err)
	}
	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("error running migrations: %w", err)
	}
	return nil
}

func setupOauth(cfg config.Config) (oauth2.Config, *oidc.Provider, error) {
	ctx := context.Background()
	var oauth2Config oauth2.Config
	var provider *oidc.Provider
	redirectUrl := fmt.Sprintf("http://localhost:%d/oauth/callback", cfg.Port)
	provider, err := oidc.NewProvider(ctx, "https://discord.com")
	if err != nil {
		return oauth2Config, provider, err
	}
	endpoints := provider.Endpoint()
	oauth2Config = oauth2.Config{
		ClientID:     cfg.Discord.ClientID,
		ClientSecret: cfg.Discord.ClientSecret,
		Endpoint:     endpoints,
		RedirectURL:  redirectUrl,
		Scopes:       []string{"identify", "guilds.join"},
	}

	return oauth2Config, provider, nil
}
