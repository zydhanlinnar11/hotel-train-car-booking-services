package hotel

import (
	"context"
	"fmt"
	"time"

	"github.com/zydhanlinnar11/hotel-train-car-booking-services/twophase/pkg/api"
	"github.com/zydhanlinnar11/hotel-train-car-booking-services/twophase/pkg/config"
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
func (s *Service) Prepare(ctx context.Context, req *api.PrepareRequest[HotelRoomReservationPayload]) (*api.PrepareResponse, error) {
	// Check if transaction already exists
	existingTransaction, err := s.repo.GetTwoPhaseTransaction(ctx, req.TransactionID)
	if err == nil && existingTransaction != nil {
		// Transaction already exists, return current status
		return &api.PrepareResponse{
			Success: existingTransaction.Status == "prepared",
			Message: fmt.Sprintf("Transaction already %s", existingTransaction.Status),
		}, nil
	}

	startDate, err := time.Parse(config.DateFormat, req.Payload.HotelRoomStartDate)
	if err != nil {
		return &api.PrepareResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to parse start date: %v", err),
		}, nil
	}

	endDate, err := time.Parse(config.DateFormat, req.Payload.HotelRoomEndDate)
	if err != nil {
		return &api.PrepareResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to parse end date: %v", err),
		}, nil
	}

	if err := s.repo.PrepareRoomReservation(
		ctx,
		req.TransactionID,
		req.Payload.HotelRoomID,
		startDate.Format(config.DateFormat),
		endDate.Format(config.DateFormat),
	); err != nil {
		return &api.PrepareResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to create transaction log: %v", err),
		}, nil
	}

	return &api.PrepareResponse{
		Success: true,
		Message: "Hotel service prepared successfully",
	}, nil
}

// Commit handles the commit phase of two-phase commit
func (s *Service) Commit(ctx context.Context, req *api.CommitRequest) (*api.CommitResponse, error) {
	if err := s.repo.CommitRoomReservation(ctx, req.TransactionID); err != nil {
		return &api.CommitResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to commit transaction: %v", err),
		}, nil
	}

	return &api.CommitResponse{
		Success: true,
		Message: "Hotel service committed successfully",
	}, nil
}

// Abort handles the abort phase of two-phase commit
func (s *Service) Abort(ctx context.Context, req *api.AbortRequest) (*api.AbortResponse, error) {
	if err := s.repo.AbortRoomReservation(ctx, req.TransactionID); err != nil {
		return &api.AbortResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to abort transaction: %v", err),
		}, nil
	}

	return &api.AbortResponse{
		Success: true,
		Message: "Hotel service aborted successfully",
	}, nil
}
