package config

import "github.com/caarlos0/env/v11"

const (
	DateFormat = "2006-01-02"
)

type Config struct {
	GoogleProjectID string `env:"GOOGLE_PROJECT_ID,required"`
}

func LoadConfig() (Config, error) {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}
