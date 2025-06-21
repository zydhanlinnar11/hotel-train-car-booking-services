package hotel

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Handler handles HTTP requests for hotel service
type Handler struct {
	service *Service
}

// NewHandler creates a new handler instance
func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

// RegisterRoutes registers all routes for hotel service
func (h *Handler) RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		// Hotel management endpoints
		api.GET("/rooms/:roomID", h.GetRoom)
		api.GET("/reservations/:orderID", h.GetReservationByOrderID)
		api.POST("/reservations", h.CreateReservation)
		api.DELETE("/reservations/:orderID", h.CancelReservation)

		// Two-phase commit endpoints
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

// GetRoom handles room retrieval
func (h *Handler) GetRoom(c *gin.Context) {
	roomID := c.Param("roomID")
	if roomID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Missing room ID",
			"message": "Room ID is required",
		})
		return
	}

	room, err := h.service.GetRoom(c.Request.Context(), roomID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Room not found",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, room)
}

// GetReservationByOrderID handles reservation retrieval by order ID
func (h *Handler) GetReservationByOrderID(c *gin.Context) {
	orderID := c.Param("orderID")
	if orderID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Missing order ID",
			"message": "Order ID is required",
		})
		return
	}

	reservation, err := h.service.GetReservationByOrderID(c.Request.Context(), orderID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Reservation not found",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, reservation)
}

// CreateReservation handles reservation creation
func (h *Handler) CreateReservation(c *gin.Context) {
	var req CreateReservationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"message": err.Error(),
		})
		return
	}

	// Validate required fields
	if req.OrderID == "" || req.HotelID == "" || req.RoomID == "" || req.UserID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Missing required fields",
			"message": "order_id, hotel_id, room_id, and user_id are required",
		})
		return
	}

	// Validate dates
	now := time.Now()
	if req.CheckInDate.Before(now) || req.CheckOutDate.Before(req.CheckInDate) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid dates",
			"message": "Check-in and check-out dates must be in the future",
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

	reservationID, err := h.service.CreateReservation(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create reservation",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"reservation_id": reservationID,
		"message":        "Reservation created successfully",
	})
}

// CancelReservation handles reservation cancellation
func (h *Handler) CancelReservation(c *gin.Context) {
	orderID := c.Param("orderID")
	if orderID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Missing order ID",
			"message": "Order ID is required",
		})
		return
	}

	err := h.service.CancelReservation(c.Request.Context(), orderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to cancel reservation",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Reservation cancelled successfully",
	})
}

// Prepare handles prepare phase requests
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

	response, err := h.service.Prepare(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to prepare transaction",
			"message": err.Error(),
		})
		return
	}

	if response.Success {
		c.JSON(http.StatusOK, response)
	} else {
		c.JSON(http.StatusBadRequest, response)
	}
}

// Commit handles commit phase requests
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

	response, err := h.service.Commit(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to commit transaction",
			"message": err.Error(),
		})
		return
	}

	if response.Success {
		c.JSON(http.StatusOK, response)
	} else {
		c.JSON(http.StatusBadRequest, response)
	}
}

// Abort handles abort phase requests
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

	response, err := h.service.Abort(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to abort transaction",
			"message": err.Error(),
		})
		return
	}

	if response.Success {
		c.JSON(http.StatusOK, response)
	} else {
		c.JSON(http.StatusBadRequest, response)
	}
}

// HealthCheck handles health check requests
func (h *Handler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"timestamp": time.Now().UTC(),
		"service":   "hotel-service",
	})
}
