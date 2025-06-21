package coordinator

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Handler handles HTTP requests for the coordinator
type Handler struct {
	service *Service
}

// NewHandler creates a new handler instance
func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

// RegisterRoutes registers all routes for the coordinator
func (h *Handler) RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		// Order creation endpoint
		api.POST("/orders", h.CreateOrder)

		// Transaction status endpoint
		api.GET("/transactions/:transactionID", h.GetTransactionStatus)

		// Two-phase commit endpoints (for participants)
		twophase := api.Group("/twophase")
		{
			twophase.POST("/prepare", h.Prepare)
			twophase.POST("/commit", h.Commit)
			twophase.POST("/abort", h.Abort)
		}

		// Health check
		api.GET("/health", h.HealthCheck)
	}
}

// CreateOrder handles order creation with two-phase commit
func (h *Handler) CreateOrder(c *gin.Context) {
	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"message": err.Error(),
		})
		return
	}

	// Validate required fields
	if req.UserID == "" || req.HotelID == "" || req.RoomID == "" ||
		req.CarID == "" || req.TrainID == "" || req.SeatID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Missing required fields",
			"message": "user_id, hotel_id, room_id, car_id, train_id, and seat_id are required",
		})
		return
	}

	// Validate dates
	now := time.Now()
	if req.CheckInDate.Before(now) || req.CheckOutDate.Before(req.CheckInDate) ||
		req.TravelDate.Before(now) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid dates",
			"message": "Check-in, check-out, and travel dates must be in the future",
		})
		return
	}

	// Validate price
	if req.TotalPrice <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid price",
			"message": "Total price must be greater than 0",
		})
		return
	}

	response, err := h.service.CreateOrder(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create order",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusAccepted, response)
}

// GetTransactionStatus handles transaction status retrieval
func (h *Handler) GetTransactionStatus(c *gin.Context) {
	transactionID := c.Param("transactionID")
	if transactionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Missing transaction ID",
			"message": "Transaction ID is required",
		})
		return
	}

	status, err := h.service.GetTransactionStatus(c.Request.Context(), transactionID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Transaction not found",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, status)
}

// Prepare handles prepare phase requests from participants
func (h *Handler) Prepare(c *gin.Context) {
	var req PrepareRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"message": err.Error(),
		})
		return
	}

	// Validate required fields
	if req.TransactionID == "" || req.OrderID == "" || req.ServiceName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Missing required fields",
			"message": "transaction_id, order_id, and service_name are required",
		})
		return
	}

	// In a real implementation, this would validate the transaction and prepare the service
	// For now, we'll just return success
	response := &PrepareResponse{
		Success: true,
		Message: "Service prepared successfully",
	}

	c.JSON(http.StatusOK, response)
}

// Commit handles commit phase requests from participants
func (h *Handler) Commit(c *gin.Context) {
	var req CommitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"message": err.Error(),
		})
		return
	}

	// Validate required fields
	if req.TransactionID == "" || req.OrderID == "" || req.ServiceName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Missing required fields",
			"message": "transaction_id, order_id, and service_name are required",
		})
		return
	}

	// In a real implementation, this would commit the transaction in the service
	// For now, we'll just return success
	response := &CommitResponse{
		Success: true,
		Message: "Service committed successfully",
	}

	c.JSON(http.StatusOK, response)
}

// Abort handles abort phase requests from participants
func (h *Handler) Abort(c *gin.Context) {
	var req AbortRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"message": err.Error(),
		})
		return
	}

	// Validate required fields
	if req.TransactionID == "" || req.OrderID == "" || req.ServiceName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Missing required fields",
			"message": "transaction_id, order_id, and service_name are required",
		})
		return
	}

	// In a real implementation, this would abort the transaction in the service
	// For now, we'll just return success
	response := &AbortResponse{
		Success: true,
		Message: "Service aborted successfully",
	}

	c.JSON(http.StatusOK, response)
}

// HealthCheck handles health check requests
func (h *Handler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"timestamp": time.Now().UTC(),
		"service":   "two-phase-commit-coordinator",
	})
}
