package hotel

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Service handles hotel business logic with two-phase commit
type Service struct {
	repo *Repository
}

// NewService creates a new hotel service
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
		Message: "Hotel service prepared successfully",
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
		Message: "Hotel service committed successfully",
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
		Message: "Hotel service aborted successfully",
	}, nil
}

// CreateReservation creates a hotel reservation (called during prepare phase)
func (s *Service) CreateReservation(ctx context.Context, req *CreateReservationRequest) (string, error) {
	// Check if room exists and is available
	room, err := s.repo.GetRoom(ctx, req.RoomID)
	if err != nil {
		return "", fmt.Errorf("room not found: %w", err)
	}

	if !room.Available {
		return "", fmt.Errorf("room is not available")
	}

	// Check room availability for the given dates
	available, err := s.repo.CheckRoomAvailability(ctx, req.RoomID, req.CheckInDate, req.CheckOutDate)
	if err != nil {
		return "", fmt.Errorf("failed to check room availability: %w", err)
	}

	if !available {
		return "", fmt.Errorf("room is not available for the specified dates")
	}

	// Create reservation
	reservation := &Reservation{
		ID:           uuid.New().String(),
		OrderID:      req.OrderID,
		HotelID:      req.HotelID,
		RoomID:       req.RoomID,
		UserID:       req.UserID,
		CheckInDate:  req.CheckInDate,
		CheckOutDate: req.CheckOutDate,
		Status:       "pending",
		TotalPrice:   req.TotalPrice,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.repo.CreateReservation(ctx, reservation); err != nil {
		return "", fmt.Errorf("failed to create reservation: %w", err)
	}

	// Update room availability
	if err := s.repo.UpdateRoomAvailability(ctx, req.RoomID, false); err != nil {
		// Try to delete the reservation we just created
		s.repo.DeleteReservation(ctx, reservation.ID)
		return "", fmt.Errorf("failed to update room availability: %w", err)
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

	// Update room availability back to available
	if err := s.repo.UpdateRoomAvailability(ctx, reservation.RoomID, true); err != nil {
		return fmt.Errorf("failed to update room availability: %w", err)
	}

	return nil
}

// GetReservationByOrderID retrieves reservation by order ID
func (s *Service) GetReservationByOrderID(ctx context.Context, orderID string) (*Reservation, error) {
	return s.repo.GetReservationByOrderID(ctx, orderID)
}

// GetRoom retrieves room information
func (s *Service) GetRoom(ctx context.Context, roomID string) (*Room, error) {
	return s.repo.GetRoom(ctx, roomID)
}
