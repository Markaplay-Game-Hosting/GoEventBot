package main

import (
	"context"
	"database/sql"
	"expvar"
	"flag"
	"fmt"
	"github.com/Markaplay-Game-Hosting/GoEventBot/internal/data"
	"github.com/Markaplay-Game-Hosting/GoEventBot/internal/vcs"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"log"
	"log/slog"
	"os"
	"runtime"
	"sync"
	"time"
)

var (
	version = vcs.Version()
)

type config struct {
	Port int    `yaml:"port"`
	Env  string `yaml:"env"`
	DB   struct {
		DSN          string `yaml:"dsn"`
		MaxOpenConns int    `yaml:"maxOpenConns"`
		MaxIdleConns int    `yaml:"maxIdleConns"`
		MaxIdleTime  string `yaml:"maxIdleTime"`
	} `yaml:"db"`
	Limiter struct {
		Enabled bool    `yaml:"enabled"`
		RPS     float64 `yaml:"rps"`
		Burst   int     `yaml:"burst"`
	} `yaml:"limiter"`
	Cors struct {
		TrustedOrigins []string `yaml:"trusted_origins"`
	} `yaml:"cors"`
	Discord struct {
		ClientID     string `yaml:"client_id"`
		ClientSecret string `yaml:"client_secret"`
	}
}

type application struct {
	config                 config
	logger                 *slog.Logger
	models                 data.Models
	scheduledEventsTracker map[uuid.UUID]data.Event
	oauth2Config           oauth2.Config
	provider               *oidc.Provider
	wg                     sync.WaitGroup
}

func main() {
	var cfg config

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	viper.SetDefault("Port", 8080)
	viper.SetDefault("Env", "development")
	viper.SetDefault("DB.DSN", "host=localhost port=5432 user=postgres password=postgres dbname=event sslmode=disable")
	viper.SetDefault("DB.MaxOpenConns", 25)
	viper.SetDefault("DB.MaxIdleConns", 25)
	viper.SetDefault("DB.MaxIdleTime", "15m")
	viper.SetDefault("Limiter.Enabled", true)
	viper.SetDefault("Limiter.RPS", 2)
	viper.SetDefault("Limiter.Burst", 4)
	viper.SetDefault("Cors.TrustedOrigins", []string{"http://localhost:3000"})
	viper.SetDefault("Discord.ClientID", "")
	viper.SetDefault("Discord.ClientSecret", "")

	if err := viper.ReadInConfig(); err != nil {
		log.Panic("Error reading config file: ", err)
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		log.Panic("Unable to decode into struct: ", err)
	}

	displayVersion := flag.Bool("version", false, "Display version and exit")

	flag.Parse()

	if *displayVersion {
		fmt.Printf("Version:\t%s\n", version)
		os.Exit(0)
	}

	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))

	db, err := openDB(cfg)
	if err != nil {
		logger.Error("Error while opening database connection: ", err)
	}
	defer db.Close()

	logger.Info("database connection pool established")

	oauth2Config, provider, err := setupOauth(cfg)
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
		config:       cfg,
		logger:       logger,
		models:       data.NewModels(db),
		oauth2Config: oauth2Config,
		provider:     provider,
	}

	err = app.serve()
	if err != nil {
		logger.Error("Error when running the app: ", "error", err)
	}
}

func openDB(cfg config) (*sql.DB, error) {
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

func setupOauth(cfg config) (oauth2.Config, *oidc.Provider, error) {
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
