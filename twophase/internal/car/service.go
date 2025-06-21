package car

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Service handles car business logic with two-phase commit
type Service struct {
	repo *Repository
}

// NewService creates a new car service
func NewService(repo *Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// Prepare handles the prepare phase of two-phase commit
func (s *Service) Prepare(ctx context.Context, req *PrepareRequest) (*PrepareResponse, error) {
	// Check if transaction already exists
	existingTransaction, err := s.repo.GetTwoPhaseTransaction(ctx, req.TransactionID)
	if err == nil && existingTransaction != nil {
		// Transaction already exists, return current status
		return &PrepareResponse{
			Success: existingTransaction.Status == "prepared",
			Message: fmt.Sprintf("Transaction already %s", existingTransaction.Status),
		}, nil
	}

	// Create two-phase transaction log
	transaction := &TwoPhaseTransaction{
		ID:            uuid.New().String(),
		TransactionID: req.TransactionID,
		OrderID:       req.OrderID,
		Status:        "prepared",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// Save transaction log
	if err := s.repo.CreateTwoPhaseTransaction(ctx, transaction); err != nil {
		return &PrepareResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to create transaction log: %v", err),
		}, nil
	}

	return &PrepareResponse{
		Success: true,
		Message: "Car service prepared successfully",
	}, nil
}

// Commit handles the commit phase of two-phase commit
func (s *Service) Commit(ctx context.Context, req *CommitRequest) (*CommitResponse, error) {
	// Get transaction log
	transaction, err := s.repo.GetTwoPhaseTransaction(ctx, req.TransactionID)
	if err != nil {
		return &CommitResponse{
			Success: false,
			Message: fmt.Sprintf("Transaction not found: %v", err),
		}, nil
	}

	// Check if already committed
	if transaction.Status == "committed" {
		return &CommitResponse{
			Success: true,
			Message: "Transaction already committed",
		}, nil
	}

	// Check if transaction is prepared
	if transaction.Status != "prepared" {
		return &CommitResponse{
			Success: false,
			Message: fmt.Sprintf("Transaction not in prepared state: %s", transaction.Status),
		}, nil
	}

	// Update transaction status to committed
	if err := s.repo.UpdateTwoPhaseTransactionStatus(ctx, req.TransactionID, "committed"); err != nil {
		return &CommitResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to update transaction status: %v", err),
		}, nil
	}

	return &CommitResponse{
		Success: true,
		Message: "Car service committed successfully",
	}, nil
}

// Abort handles the abort phase of two-phase commit
func (s *Service) Abort(ctx context.Context, req *AbortRequest) (*AbortResponse, error) {
	// Get transaction log
	transaction, err := s.repo.GetTwoPhaseTransaction(ctx, req.TransactionID)
	if err != nil {
		return &AbortResponse{
			Success: false,
			Message: fmt.Sprintf("Transaction not found: %v", err),
		}, nil
	}

	// Check if already aborted
	if transaction.Status == "aborted" {
		return &AbortResponse{
			Success: true,
			Message: "Transaction already aborted",
		}, nil
	}

	// If there's a reservation, delete it
	if transaction.ReservationID != "" {
		if err := s.repo.DeleteReservation(ctx, transaction.ReservationID); err != nil {
			// Log error but continue with abort
			fmt.Printf("Failed to delete reservation during abort: %v\n", err)
		}
	}

	// Update transaction status to aborted
	if err := s.repo.UpdateTwoPhaseTransactionStatus(ctx, req.TransactionID, "aborted"); err != nil {
		return &AbortResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to update transaction status: %v", err),
		}, nil
	}

	return &AbortResponse{
		Success: true,
		Message: "Car service aborted successfully",
	}, nil
}

// CreateReservation creates a car reservation (called during prepare phase)
func (s *Service) CreateReservation(ctx context.Context, req *CreateReservationRequest) (string, error) {
	// Check if car exists and is available
	car, err := s.repo.GetCar(ctx, req.CarID)
	if err != nil {
		return "", fmt.Errorf("car not found: %w", err)
	}

	if !car.Available {
		return "", fmt.Errorf("car is not available")
	}

	// Check car availability for the given dates
	available, err := s.repo.CheckCarAvailability(ctx, req.CarID, req.StartDate, req.EndDate)
	if err != nil {
		return "", fmt.Errorf("failed to check car availability: %w", err)
	}

	if !available {
		return "", fmt.Errorf("car is not available for the specified dates")
	}

	// Create reservation
	reservation := &Reservation{
		ID:         uuid.New().String(),
		OrderID:    req.OrderID,
		CarID:      req.CarID,
		UserID:     req.UserID,
		StartDate:  req.StartDate,
		EndDate:    req.EndDate,
		Status:     "pending",
		TotalPrice: req.TotalPrice,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := s.repo.CreateReservation(ctx, reservation); err != nil {
		return "", fmt.Errorf("failed to create reservation: %w", err)
	}

	// Update car availability
	if err := s.repo.UpdateCarAvailability(ctx, req.CarID, false); err != nil {
		// Try to delete the reservation we just created
		s.repo.DeleteReservation(ctx, reservation.ID)
		return "", fmt.Errorf("failed to update car availability: %w", err)
	}

	return reservation.ID, nil
}

// CancelReservation cancels a reservation by order ID
func (s *Service) CancelReservation(ctx context.Context, orderID string) error {
	// Get reservation by order ID
	reservation, err := s.repo.GetReservationByOrderID(ctx, orderID)
	if err != nil {
		return fmt.Errorf("reservation not found: %w", err)
	}

	// Update reservation status to cancelled
	if err := s.repo.UpdateReservationStatus(ctx, reservation.ID, "cancelled"); err != nil {
		return fmt.Errorf("failed to update reservation status: %w", err)
	}

	// Update car availability back to available
	if err := s.repo.UpdateCarAvailability(ctx, reservation.CarID, true); err != nil {
		return fmt.Errorf("failed to update car availability: %w", err)
	}

	return nil
}

// GetReservationByOrderID retrieves reservation by order ID
func (s *Service) GetReservationByOrderID(ctx context.Context, orderID string) (*Reservation, error) {
	return s.repo.GetReservationByOrderID(ctx, orderID)
}

// GetCar retrieves car information
func (s *Service) GetCar(ctx context.Context, carID string) (*Car, error) {
	return s.repo.GetCar(ctx, carID)
}
