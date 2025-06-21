package config

import "github.com/caarlos0/env/v11"

const (
	DateFormat = "02-01-2006"
)

type Config struct {
	RabbitMQURL     string `env:"RABBITMQ_URL,required"`
	GoogleProjectID string `env:"GOOGLE_PROJECT_ID,required"`
	Port            string `env:"PORT" envDefault:"8080"`

	OrderQueueName string `env:"ORDER_QUEUE_NAME" envDefault:"order_service_queue"`
	HotelQueueName string `env:"HOTEL_QUEUE_NAME" envDefault:"hotel_service_queue"`
	CarQueueName   string `env:"CAR_QUEUE_NAME" envDefault:"car_service_queue"`
	TrainQueueName string `env:"TRAIN_QUEUE_NAME" envDefault:"train_service_queue"`
}

func LoadConfig() (Config, error) {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}
