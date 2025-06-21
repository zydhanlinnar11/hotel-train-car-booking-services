package coordinator

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// Service handles the two-phase commit coordination logic
type Service struct {
	repo   *Repository
	config *Config
	client *http.Client
}

// NewService creates a new coordinator service
func NewService(repo *Repository, config *Config) *Service {
	return &Service{
		repo:   repo,
		config: config,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// CreateOrder initiates a two-phase commit transaction for order creation
func (s *Service) CreateOrder(ctx context.Context, req *CreateOrderRequest) (*OrderResponse, error) {
	// Generate transaction ID
	transactionID := uuid.New().String()
	orderID := uuid.New().String()

	// Create transaction log
	log := &TransactionLog{
		ID:         transactionID,
		OrderID:    orderID,
		Status:     StatusInitiated,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		TimeoutAt:  time.Now().Add(s.config.TransactionTimeout),
		RetryCount: 0,
		MaxRetries: s.config.MaxRetries,
		Participants: []Participant{
			{ServiceName: "hotel", ServiceURL: s.config.Services["hotel"], Status: "pending"},
			{ServiceName: "car", ServiceURL: s.config.Services["car"], Status: "pending"},
			{ServiceName: "train", ServiceURL: s.config.Services["train"], Status: "pending"},
		},
	}

	// Save transaction log
	if err := s.repo.CreateTransactionLog(ctx, log); err != nil {
		return nil, fmt.Errorf("failed to create transaction log: %w", err)
	}

	// Start two-phase commit in background
	go s.executeTwoPhaseCommit(context.Background(), transactionID, req)

	return &OrderResponse{
		OrderID:       orderID,
		TransactionID: transactionID,
		Status:        StatusInitiated,
		Message:       "Transaction initiated successfully",
	}, nil
}

// executeTwoPhaseCommit executes the two-phase commit protocol
func (s *Service) executeTwoPhaseCommit(ctx context.Context, transactionID string, req *CreateOrderRequest) {
	// Phase 1: Prepare
	if !s.preparePhase(ctx, transactionID, req) {
		s.abortTransaction(ctx, transactionID, "Prepare phase failed")
		return
	}

	// Phase 2: Commit
	if !s.commitPhase(ctx, transactionID) {
		s.rollbackTransaction(ctx, transactionID, "Commit phase failed")
		return
	}

	// Mark transaction as committed
	s.finalizeTransaction(ctx, transactionID, StatusCommitted, "")
}

// preparePhase executes the prepare phase
func (s *Service) preparePhase(ctx context.Context, transactionID string, req *CreateOrderRequest) bool {
	log, err := s.repo.GetTransactionLog(ctx, transactionID)
	if err != nil {
		return false
	}

	// Update status to prepared
	log.Status = StatusPrepared
	log.UpdatedAt = time.Now()
	if err := s.repo.UpdateTransactionLog(ctx, log); err != nil {
		return false
	}

	// Send prepare requests to all participants
	prepareReq := &PrepareRequest{
		TransactionID: transactionID,
		OrderID:       log.OrderID,
	}

	allPrepared := true
	for _, participant := range log.Participants {
		if !s.sendPrepareRequest(ctx, transactionID, participant.ServiceName, prepareReq) {
			allPrepared = false
			break
		}
	}

	return allPrepared
}

// commitPhase executes the commit phase
func (s *Service) commitPhase(ctx context.Context, transactionID string) bool {
	log, err := s.repo.GetTransactionLog(ctx, transactionID)
	if err != nil {
		return false
	}

	// Send commit requests to all participants
	commitReq := &CommitRequest{
		TransactionID: transactionID,
		OrderID:       log.OrderID,
	}

	allCommitted := true
	for _, participant := range log.Participants {
		if !s.sendCommitRequest(ctx, transactionID, participant.ServiceName, commitReq) {
			allCommitted = false
			break
		}
	}

	return allCommitted
}

// sendPrepareRequest sends prepare request to a participant with retry logic
func (s *Service) sendPrepareRequest(ctx context.Context, transactionID, serviceName string, req *PrepareRequest) bool {
	serviceURL := s.config.Services[serviceName]
	url := fmt.Sprintf("%s/api/twophase/prepare", serviceURL)

	return s.sendRequestWithRetry(ctx, transactionID, serviceName, url, req, "prepare")
}

// sendCommitRequest sends commit request to a participant with retry logic
func (s *Service) sendCommitRequest(ctx context.Context, transactionID, serviceName string, req *CommitRequest) bool {
	serviceURL := s.config.Services[serviceName]
	url := fmt.Sprintf("%s/api/twophase/commit", serviceURL)

	return s.sendRequestWithRetry(ctx, transactionID, serviceName, url, req, "commit")
}

// sendRequestWithRetry sends HTTP request with exponential backoff retry
func (s *Service) sendRequestWithRetry(ctx context.Context, transactionID, serviceName, url string, payload interface{}, operation string) bool {
	maxRetries := s.config.MaxRetries
	baseDelay := s.config.RetryDelay

	for attempt := 0; attempt <= maxRetries; attempt++ {
		success := s.sendSingleRequest(ctx, transactionID, serviceName, url, payload, operation)
		if success {
			return true
		}

		if attempt < maxRetries {
			// Exponential backoff
			delay := time.Duration(float64(baseDelay) * math.Pow(2, float64(attempt)))
			time.Sleep(delay)
		}
	}

	// Update participant status to failed
	s.repo.UpdateParticipantStatus(ctx, transactionID, serviceName, "failed", fmt.Sprintf("%s operation failed after %d retries", operation, maxRetries))
	return false
}

// sendSingleRequest sends a single HTTP request
func (s *Service) sendSingleRequest(ctx context.Context, transactionID, serviceName, url string, payload interface{}, operation string) bool {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return false
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return false
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		// Update participant status
		status := "prepared"
		if operation == "commit" {
			status = "committed"
		}
		s.repo.UpdateParticipantStatus(ctx, transactionID, serviceName, status, "")
		return true
	}

	return false
}

// abortTransaction aborts the transaction
func (s *Service) abortTransaction(ctx context.Context, transactionID, reason string) {
	log, err := s.repo.GetTransactionLog(ctx, transactionID)
	if err != nil {
		return
	}

	// Send abort requests to all participants
	abortReq := &AbortRequest{
		TransactionID: transactionID,
		OrderID:       log.OrderID,
		Reason:        reason,
	}

	for _, participant := range log.Participants {
		s.sendAbortRequest(ctx, transactionID, participant.ServiceName, abortReq)
	}

	s.finalizeTransaction(ctx, transactionID, StatusAborted, reason)
}

// rollbackTransaction rolls back the transaction
func (s *Service) rollbackTransaction(ctx context.Context, transactionID, reason string) {
	s.finalizeTransaction(ctx, transactionID, StatusRolledBack, reason)
}

// sendAbortRequest sends abort request to a participant
func (s *Service) sendAbortRequest(ctx context.Context, transactionID, serviceName string, req *AbortRequest) {
	serviceURL := s.config.Services[serviceName]
	url := fmt.Sprintf("%s/api/twophase/abort", serviceURL)

	jsonData, err := json.Marshal(req)
	if err != nil {
		return
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return
	}

	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(httpReq)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		s.repo.UpdateParticipantStatus(ctx, transactionID, serviceName, "aborted", "")
	}
}

// finalizeTransaction finalizes the transaction status
func (s *Service) finalizeTransaction(ctx context.Context, transactionID string, status TransactionStatus, reason string) {
	log, err := s.repo.GetTransactionLog(ctx, transactionID)
	if err != nil {
		return
	}

	log.Status = status
	log.UpdatedAt = time.Now()
	log.FailureReason = reason

	if status == StatusCommitted {
		now := time.Now()
		log.CommitTimestamp = &now
	}

	s.repo.UpdateTransactionLog(ctx, log)
}

// GetTransactionStatus retrieves the status of a transaction
func (s *Service) GetTransactionStatus(ctx context.Context, transactionID string) (*TransactionStatusResponse, error) {
	log, err := s.repo.GetTransactionLog(ctx, transactionID)
	if err != nil {
		return nil, err
	}

	return &TransactionStatusResponse{
		TransactionID: log.ID,
		OrderID:       log.OrderID,
		Status:        log.Status,
		Participants:  log.Participants,
		CreatedAt:     log.CreatedAt,
		UpdatedAt:     log.UpdatedAt,
		TimeoutAt:     log.TimeoutAt,
		RetryCount:    log.RetryCount,
		FailureReason: log.FailureReason,
	}, nil
}

// CleanupTimedOutTransactions cleans up timed out transactions
func (s *Service) CleanupTimedOutTransactions(ctx context.Context) error {
	timedOutLogs, err := s.repo.GetTimedOutTransactions(ctx)
	if err != nil {
		return err
	}

	for _, log := range timedOutLogs {
		log.Status = StatusTimedOut
		log.UpdatedAt = time.Now()
		log.FailureReason = "Transaction timed out"

		if err := s.repo.UpdateTransactionLog(ctx, log); err != nil {
			continue
		}

		// Send abort requests to participants
		abortReq := &AbortRequest{
			TransactionID: log.ID,
			OrderID:       log.OrderID,
			Reason:        "Transaction timed out",
		}

		for _, participant := range log.Participants {
			s.sendAbortRequest(ctx, log.ID, participant.ServiceName, abortReq)
		}
	}

	return nil
}
