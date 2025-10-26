package model

// ProbeExecutionRequest represents a request to manage probe execution
type ProbeExecutionRequest struct {
	Action       string        `json:"action" binding:"required"`
	ResourceName string        `json:"resource_name" binding:"required"`
	Namespace    string        `json:"namespace" binding:"required"`
	Policy       PolicyDetails `json:"policy" binding:"required"`
}

// PolicyDetails contains the policy configuration for probe execution
type PolicyDetails struct {
	Name        string            `json:"name"`
	Type        string            `json:"type"`
	Target      string            `json:"target"`
	Parameters  map[string]string `json:"parameters,omitempty"`
	Enabled     bool              `json:"enabled"`
	Description string            `json:"description,omitempty"`
}

// ProbeExecutionStatus represents the current status of a probe execution
type ProbeExecutionStatus struct {
	ResourceName string `json:"resource_name"`
	Namespace    string `json:"namespace"`
	Status       string `json:"status"`
	Message      string `json:"message,omitempty"`
	Policy       string `json:"policy,omitempty"`
}

// APIResponse represents a generic API response
type APIResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
