package hotel

import (
	"net/http"

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
	// Two-phase commit endpoints
	twophase := r.Group("/twophase")
	{
		twophase.POST("/prepare", h.Prepare)
		twophase.POST("/commit", h.Commit)
		twophase.POST("/abort", h.Abort)
	}

	// Health check
	// r.GET("/health", h.HealthCheck)
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
