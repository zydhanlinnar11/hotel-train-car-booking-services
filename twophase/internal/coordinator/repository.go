package coordinator

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Repository handles Firestore operations for transaction logs
type Repository struct {
	client *firestore.Client
}

// NewRepository creates a new repository instance
func NewRepository(client *firestore.Client) *Repository {
	return &Repository{
		client: client,
	}
}

// CreateTransactionLog creates a new transaction log entry
func (r *Repository) CreateTransactionLog(ctx context.Context, log *TransactionLog) error {
	collection := r.client.Collection("twophase_transactions")
	doc := collection.Doc(log.ID)

	_, err := doc.Create(ctx, log)
	if err != nil {
		return fmt.Errorf("failed to create transaction log: %w", err)
	}

	return nil
}

// GetTransactionLog retrieves a transaction log by ID
func (r *Repository) GetTransactionLog(ctx context.Context, transactionID string) (*TransactionLog, error) {
	collection := r.client.Collection("twophase_transactions")
	doc := collection.Doc(transactionID)

	docSnap, err := doc.Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, fmt.Errorf("transaction log not found: %s", transactionID)
		}
		return nil, fmt.Errorf("failed to get transaction log: %w", err)
	}

	var log TransactionLog
	if err := docSnap.DataTo(&log); err != nil {
		return nil, fmt.Errorf("failed to unmarshal transaction log: %w", err)
	}

	return &log, nil
}

// UpdateTransactionLog updates an existing transaction log
func (r *Repository) UpdateTransactionLog(ctx context.Context, log *TransactionLog) error {
	collection := r.client.Collection("twophase_transactions")
	doc := collection.Doc(log.ID)

	log.UpdatedAt = time.Now()

	_, err := doc.Set(ctx, log)
	if err != nil {
		return fmt.Errorf("failed to update transaction log: %w", err)
	}

	return nil
}

// UpdateParticipantStatus updates the status of a specific participant
func (r *Repository) UpdateParticipantStatus(ctx context.Context, transactionID, serviceName, status, errorMsg string) error {
	collection := r.client.Collection("twophase_transactions")
	doc := collection.Doc(transactionID)

	// Get current transaction log
	docSnap, err := doc.Get(ctx)
	if err != nil {
		return fmt.Errorf("failed to get transaction log: %w", err)
	}

	var log TransactionLog
	if err := docSnap.DataTo(&log); err != nil {
		return fmt.Errorf("failed to unmarshal transaction log: %w", err)
	}

	// Update participant status
	for i, participant := range log.Participants {
		if participant.ServiceName == serviceName {
			log.Participants[i].Status = status
			log.Participants[i].Error = errorMsg
			if status == "failed" {
				log.Participants[i].RetryCount++
			}
			break
		}
	}

	log.UpdatedAt = time.Now()

	_, err = doc.Set(ctx, log)
	if err != nil {
		return fmt.Errorf("failed to update participant status: %w", err)
	}

	return nil
}

// GetTimedOutTransactions retrieves transactions that have timed out
func (r *Repository) GetTimedOutTransactions(ctx context.Context) ([]*TransactionLog, error) {
	collection := r.client.Collection("twophase_transactions")
	now := time.Now()

	query := collection.Where("timeout_at", "<=", now).
		Where("status", "in", []string{string(StatusInitiated), string(StatusPrepared)})

	iter := query.Documents(ctx)
	defer iter.Stop()

	var logs []*TransactionLog
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to iterate transaction logs: %w", err)
		}

		var log TransactionLog
		if err := doc.DataTo(&log); err != nil {
			return nil, fmt.Errorf("failed to unmarshal transaction log: %w", err)
		}

		logs = append(logs, &log)
	}

	return logs, nil
}

// GetPendingTransactions retrieves transactions that are still pending
func (r *Repository) GetPendingTransactions(ctx context.Context) ([]*TransactionLog, error) {
	collection := r.client.Collection("twophase_transactions")

	query := collection.Where("status", "in", []string{string(StatusInitiated), string(StatusPrepared)})

	iter := query.Documents(ctx)
	defer iter.Stop()

	var logs []*TransactionLog
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to iterate transaction logs: %w", err)
		}

		var log TransactionLog
		if err := doc.DataTo(&log); err != nil {
			return nil, fmt.Errorf("failed to unmarshal transaction log: %w", err)
		}

		logs = append(logs, &log)
	}

	return logs, nil
}

// DeleteTransactionLog deletes a transaction log (for cleanup purposes)
func (r *Repository) DeleteTransactionLog(ctx context.Context, transactionID string) error {
	collection := r.client.Collection("twophase_transactions")
	doc := collection.Doc(transactionID)

	_, err := doc.Delete(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete transaction log: %w", err)
	}

	return nil
}

// GetTransactionLogsByOrderID retrieves transaction logs by order ID
func (r *Repository) GetTransactionLogsByOrderID(ctx context.Context, orderID string) ([]*TransactionLog, error) {
	collection := r.client.Collection("twophase_transactions")

	query := collection.Where("order_id", "==", orderID)

	iter := query.Documents(ctx)
	defer iter.Stop()

	var logs []*TransactionLog
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to iterate transaction logs: %w", err)
		}

		var log TransactionLog
		if err := doc.DataTo(&log); err != nil {
			return nil, fmt.Errorf("failed to unmarshal transaction log: %w", err)
		}

		logs = append(logs, &log)
	}

	return logs, nil
}
