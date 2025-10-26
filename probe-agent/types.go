package main

// PolicyEventRequest represents an event request from the operator
type PolicyEventRequest struct {
	Action       string        `json:"action" binding:"required"`
	ResourceName string        `json:"resource_name" binding:"required"`
	Namespace    string        `json:"namespace" binding:"required"`
	Policy       PolicyDetails `json:"policy" binding:"required"`
}

// PolicyDetails contains the policy configuration
type PolicyDetails struct {
	Name        string            `json:"name"`
	Type        string            `json:"type"`
	Target      string            `json:"target"`
	Parameters  map[string]string `json:"parameters,omitempty"`
	Enabled     bool              `json:"enabled"`
	Description string            `json:"description,omitempty"`
}

// EventResponse represents the response sent back to the operator
type EventResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}
