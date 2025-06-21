package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/zydhanlinnar11/hotel-train-car-booking-services/twophase/internal/coordinator"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Initialize Firestore client
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, os.Getenv("GOOGLE_CLOUD_PROJECT"))
	if err != nil {
		log.Fatalf("Failed to create Firestore client: %v", err)
	}
	defer client.Close()

	// Initialize repository
	repo := coordinator.NewRepository(client)

	// Initialize configuration
	config := coordinator.DefaultConfig()

	// Override with environment variables if provided
	if timeout := os.Getenv("TRANSACTION_TIMEOUT"); timeout != "" {
		if duration, err := time.ParseDuration(timeout); err == nil {
			config.TransactionTimeout = duration
		}
	}

	if maxRetries := os.Getenv("MAX_RETRIES"); maxRetries != "" {
		if retries, err := time.ParseDuration(maxRetries); err == nil {
			config.MaxRetries = int(retries)
		}
	}

	if retryDelay := os.Getenv("RETRY_DELAY"); retryDelay != "" {
		if delay, err := time.ParseDuration(retryDelay); err == nil {
			config.RetryDelay = delay
		}
	}

	// Override service URLs if provided
	if hotelURL := os.Getenv("HOTEL_SERVICE_URL"); hotelURL != "" {
		config.Services["hotel"] = hotelURL
	}
	if carURL := os.Getenv("CAR_SERVICE_URL"); carURL != "" {
		config.Services["car"] = carURL
	}
	if trainURL := os.Getenv("TRAIN_SERVICE_URL"); trainURL != "" {
		config.Services["train"] = trainURL
	}

	// Initialize service
	service := coordinator.NewService(repo, config)

	// Initialize handler
	handler := coordinator.NewHandler(service)

	// Initialize Gin router
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// Add CORS middleware
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Register routes
	handler.RegisterRoutes(r)

	// Start cleanup goroutine
	go func() {
		ticker := time.NewTicker(1 * time.Minute) // Check every minute
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if err := service.CleanupTimedOutTransactions(ctx); err != nil {
					log.Printf("Failed to cleanup timed out transactions: %v", err)
				}
			}
		}
	}()

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Starting coordinator service on port %s", port)
		if err := r.Run(":" + port); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down coordinator service...")

	// Give outstanding requests a chance to complete
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	log.Println("Coordinator service stopped")
}
