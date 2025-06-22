package api

// PrepareResponse represents prepare phase response
type PrepareResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// CommitResponse represents commit phase response
type CommitResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// AbortResponse represents abort phase response
type AbortResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
