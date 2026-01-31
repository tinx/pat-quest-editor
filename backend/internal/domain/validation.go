package domain

import "fmt"

// ValidationError represents a single validation issue.
type ValidationError struct {
	NodeID  *int   `json:"nodeId,omitempty"`
	Field   string `json:"field,omitempty"`
	Message string `json:"message"`
}

func (e ValidationError) Error() string {
	if e.NodeID != nil {
		return fmt.Sprintf("Node %d: %s", *e.NodeID, e.Message)
	}
	return e.Message
}

// ValidationResult contains all validation errors for a quest.
type ValidationResult struct {
	Valid  bool              `json:"valid"`
	Errors []ValidationError `json:"errors,omitempty"`
}

// NewValidationResult creates an empty valid result.
func NewValidationResult() *ValidationResult {
	return &ValidationResult{Valid: true, Errors: []ValidationError{}}
}

// AddError adds an error and marks result as invalid.
func (r *ValidationResult) AddError(err ValidationError) {
	r.Valid = false
	r.Errors = append(r.Errors, err)
}

// AddNodeError adds an error associated with a specific node.
func (r *ValidationResult) AddNodeError(nodeID int, message string) {
	r.AddError(ValidationError{NodeID: &nodeID, Message: message})
}

// AddGlobalError adds an error not associated with a specific node.
func (r *ValidationResult) AddGlobalError(message string) {
	r.AddError(ValidationError{Message: message})
}
