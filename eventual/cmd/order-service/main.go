package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/rabbitmq/amqp091-go"
	"github.com/zydhanlinnar11/hotel-train-car-booking-services/eventual/internal/order"
	"github.com/zydhanlinnar11/hotel-train-car-booking-services/eventual/pkg/config"
	"github.com/zydhanlinnar11/hotel-train-car-booking-services/eventual/pkg/event"
	"github.com/zydhanlinnar11/hotel-train-car-booking-services/eventual/pkg/messagebus"
)

func main() {
	// Create a cancellable context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log.Println("Starting order service")
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

	orderRepo := order.NewFirestoreRepository(client)
	orderService := order.NewService(orderRepo, publisher)
	orderHandler := order.NewHandler(orderService)

	subscriber := messagebus.NewRabbitmqSubscriber(conn)
	subscriber.Subscribe(ctx, string(event.CommandReserveRoom), "order-service", func(e event.Message) {
		orderService.ProcessSagaEvent(ctx, e)
	})
	subscriber.Subscribe(ctx, string(event.CommandReserveCar), "order-service", func(e event.Message) {
		orderService.ProcessSagaEvent(ctx, e)
	})
	subscriber.Subscribe(ctx, string(event.CommandReserveSeat), "order-service", func(e event.Message) {
		orderService.ProcessSagaEvent(ctx, e)
	})

	// Start HTTP server
	router := gin.Default()
	router.POST("/orders", orderHandler.CreateOrder)

	// Create HTTP server with proper shutdown handling
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	log.Println("Order service started")

	// Start HTTP server in a goroutine
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start HTTP server: %v", err)
		}
	}()

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

	// Graceful shutdown with timeout
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	// Shutdown HTTP server
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("HTTP server shutdown error: %v", err)
	}

	// Note: Subscriber goroutines will be terminated when context is cancelled
	// and connection is closed. The current implementation doesn't provide
	// explicit unsubscribe functionality.

	log.Println("Order service stopped gracefully")
}
