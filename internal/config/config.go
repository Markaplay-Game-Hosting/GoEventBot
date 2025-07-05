package config

import (
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	Port int
	Env  string
	DB   struct {
		DSN          string
		MaxOpenConns int
		MaxIdleConns int
		MaxIdleTime  string
	}
	Limiter struct {
		Enabled bool
		RPS     float64
		Burst   int
	}
	Cors struct {
		TrustedOrigins []string
	}
	Discord struct {
		ClientID     string
		ClientSecret string
		Token        string
	}
	Security struct {
		Secret      string
		SecretKey   string
		SecretBlock string
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
	viper.SetDefault("Discord.Token", "")
	viper.SetDefault("Security.Secret", "")
	viper.SetDefault("Security.SecretBlock", "")
	viper.SetDefault("Security.SecretKey", "")

	if err := viper.ReadInConfig(); err != nil {
		log.Panic("Error reading config file: ", err)
		return nil, err
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		log.Panic("Unable to decode into struct: ", err)
	}

	return &cfg, nil
}
