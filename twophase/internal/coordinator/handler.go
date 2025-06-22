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
	// Order creation endpoint
	r.POST("/orders", h.CreateOrder)

	// Transaction status endpoint
	r.GET("/transactions/:transactionID", h.GetTransactionStatus)

	// Health check
	r.GET("/health", h.HealthCheck)
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
	if req.UserID == "" || req.HotelRoomID == "" || req.CarID == "" ||
		req.TrainSeatID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Missing required fields",
			"message": "user_id, hotel_id, room_id, car_id, train_id, and seat_id are required",
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

// HealthCheck handles health check requests
func (h *Handler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"timestamp": time.Now().UTC(),
		"service":   "two-phase-commit-coordinator",
	})
}
