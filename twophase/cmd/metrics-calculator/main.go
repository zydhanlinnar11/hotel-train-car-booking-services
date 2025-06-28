package main

import (
	"context"
	"encoding/csv"
	"log"
	"os"
	"strconv"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/joho/godotenv"
	"github.com/zydhanlinnar11/hotel-train-car-booking-services/twophase/internal/coordinator"
	"github.com/zydhanlinnar11/hotel-train-car-booking-services/twophase/pkg/config"
	"google.golang.org/api/iterator"
)

func main() {
	log.Println("Starting metrics calculator...")

	// Validasi argument untuk nama file export
	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s <output_filename.csv>", os.Args[0])
	}

	outputFilename := os.Args[1]
	if outputFilename == "" {
		log.Fatalf("Output filename cannot be empty")
	}

	if err := godotenv.Load(); err != nil {
		log.Printf("Error loading .env file: %v, using system environment variables", err)
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	ctx := context.Background()
	client, err := firestore.NewClient(ctx, cfg.GoogleProjectID)
	if err != nil {
		log.Fatalf("Failed to create Firestore client: %v", err)
	}
	defer client.Close()

	transactionLogs := make([]coordinator.TransactionLog, 0)

	iter := client.Collection("twophase_transactions").Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to iterate transaction logs: %v", err)
		}

		var tl coordinator.TransactionLog
		if err := doc.DataTo(&tl); err != nil {
			log.Fatalf("Failed to convert document to transaction log: %v", err)
		}
		transactionLogs = append(transactionLogs, tl)
	}

	// Export to CSV dengan nama file yang dikustomisasi
	if err := exportTransactionLogsToCSV(transactionLogs, outputFilename); err != nil {
		log.Fatalf("Failed to export transaction logs to CSV: %v", err)
	}

	log.Printf("Successfully exported %d transaction logs to %s", len(transactionLogs), outputFilename)
}

func exportTransactionLogsToCSV(transactionLogs []coordinator.TransactionLog, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	headers := []string{
		"ID",
		"OrderID",
		"Status",
		"Participants",
		"ParticipantsDoneAt",
		"DoneAt",
		"CreatedAt",
		"UpdatedAt",
		"TimeoutAt",
		"RetryCount",
		"MaxRetries",
		"LastRetryAt",
		"FailureReason",
		"CommitTimestamp",
	}

	if err := writer.Write(headers); err != nil {
		return err
	}

	// Write data rows
	for _, tl := range transactionLogs {
		// Convert participants to string representation
		participantsStr := ""
		participantsDoneAtStr := ""
		for i, participant := range tl.Participants {
			if i > 0 {
				participantsStr += ";"
				participantsDoneAtStr += ";"
			}
			participantsStr += participant.ServiceName + ":" + participant.Status
			participantsDoneAtStr += participant.ServiceName + ":" + formatTimePtr(participant.DoneAt)
		}

		row := []string{
			tl.ID,
			tl.OrderID,
			string(tl.Status),
			participantsStr,
			participantsDoneAtStr,
			formatTimePtr(tl.DoneAt),
			strconv.FormatInt(formatTime(tl.CreatedAt), 10),
			strconv.FormatInt(formatTime(tl.UpdatedAt), 10),
			strconv.FormatInt(formatTime(tl.TimeoutAt), 10),
			strconv.Itoa(tl.RetryCount),
			strconv.Itoa(tl.MaxRetries),
			formatTimePtr(tl.LastRetryAt),
			tl.FailureReason,
			formatTimePtr(tl.CommitTimestamp),
		}

		if err := writer.Write(row); err != nil {
			return err
		}
	}

	return nil
}

func formatTime(t time.Time) int64 {
	return t.UnixMilli()
}

func formatTimePtr(t *time.Time) string {
	if t == nil {
		return ""
	}
	return strconv.FormatInt(t.UnixMilli(), 10)
}
