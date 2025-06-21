package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"cloud.google.com/go/firestore"
	"github.com/joho/godotenv"
	"github.com/rabbitmq/amqp091-go"
	"github.com/zydhanlinnar11/hotel-train-car-booking-services/eventual/internal/train"
	"github.com/zydhanlinnar11/hotel-train-car-booking-services/eventual/pkg/config"
	"github.com/zydhanlinnar11/hotel-train-car-booking-services/eventual/pkg/event"
	"github.com/zydhanlinnar11/hotel-train-car-booking-services/eventual/pkg/messagebus"
)

func main() {
	// Create a cancellable context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log.Println("Starting train service")
	if err := godotenv.Load(); err != nil {
		log.Printf("Error loading .env file: %v, using system environment variables", err)
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	conn, err := amqp091.Dial(cfg.RabbitMQURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	client, err := firestore.NewClient(ctx, cfg.GoogleProjectID)
	if err != nil {
		log.Fatalf("Failed to create Firestore client: %v", err)
	}
	defer client.Close()

	publisher := messagebus.NewRabbitmqPublisher(conn)

	trainRepo := train.NewFirestoreRepository(client)
	trainService := train.NewService(trainRepo, publisher)

	subscriber := messagebus.NewRabbitmqSubscriber(conn)
	subscriber.Subscribe(ctx, "", cfg.TrainQueueName, func(e event.Message) {
		trainService.ProcessSagaEvent(ctx, e)
	})

	log.Println("Train service started")

	// Wait for interrupt signal to gracefully shutdown the server
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-ctx.Done():
		log.Println("Context done, shutting down")
	case <-signals:
		log.Println("Received shutdown signal, shutting down")
	}

	// Cancel context to stop all operations
	cancel()

	// Note: Subscriber goroutines will be terminated when context is cancelled
	// and connection is closed. The current implementation doesn't provide
	// explicit unsubscribe functionality.

	log.Println("Train service stopped gracefully")
}
