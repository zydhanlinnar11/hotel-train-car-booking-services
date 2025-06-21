package config

import "github.com/caarlos0/env/v11"

const (
	DateFormat = "02-01-2006"
)

type Config struct {
	RabbitMQURL     string `env:"RABBITMQ_URL,required"`
	GoogleProjectID string `env:"GOOGLE_PROJECT_ID,required"`
	Port            string `env:"PORT" envDefault:"8080"`
}

func LoadConfig() (Config, error) {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}
