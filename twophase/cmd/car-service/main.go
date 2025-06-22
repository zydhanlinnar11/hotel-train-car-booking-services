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
	"github.com/zydhanlinnar11/hotel-train-car-booking-services/twophase/internal/car"
	"github.com/zydhanlinnar11/hotel-train-car-booking-services/twophase/pkg/config"
)

const (
	Port = "8082"
)

func main() {
	// Create a cancellable context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log.Println("Starting car service")
	if err := godotenv.Load(); err != nil {
		log.Printf("Error loading .env file: %v, using system environment variables", err)
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	client, err := firestore.NewClient(ctx, cfg.GoogleProjectID)
	if err != nil {
		log.Fatalf("Failed to create Firestore client: %v", err)
	}
	defer client.Close()

	carRepo := car.NewRepository(client)
	carService := car.NewService(carRepo)
	carHandler := car.NewHandler(carService)

	// Start HTTP server
	router := gin.Default()
	carHandler.RegisterRoutes(router)

	// Create HTTP server with proper shutdown handling
	srv := &http.Server{
		Addr:    ":" + Port,
		Handler: router,
	}

	log.Println("Car service started at port", Port)

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

	log.Println("Car service stopped gracefully")
}
