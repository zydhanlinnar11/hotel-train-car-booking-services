package api

// PrepareRequest represents prepare phase request
type PrepareRequest[T any] struct {
	TransactionID string `json:"transaction_id"`
	OrderID       string `json:"order_id"`
	Payload       T      `json:"payload"`
}

// CommitRequest represents commit phase request
type CommitRequest struct {
	TransactionID string `json:"transaction_id" binding:"required"`
}

// AbortRequest represents abort phase request
type AbortRequest struct {
	TransactionID string `json:"transaction_id" binding:"required"`
}
