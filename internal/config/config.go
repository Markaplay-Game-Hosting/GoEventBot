package config

import (
	"github.com/spf13/viper"
	"log"
)

type Config struct {
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

func Load() (*Config, error) {
	var cfg Config

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
		return nil, err
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		log.Panic("Unable to decode into struct: ", err)
	}

	return &cfg, nil
}
