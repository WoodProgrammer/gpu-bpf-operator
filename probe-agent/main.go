package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("Starting GPU EBPF Probe policy agent v0.0.1")

	// Create Gin router with default middleware
	router := gin.Default()

	// Health check endpoint
	router.GET("/health", healthCheck)

	// Policy event handler endpoint
	router.POST("/events", handlePolicyEvent)

	// Start server
	port := ":8080"
	log.Printf("Probe agent listening on %s", port)
	if err := router.Run(port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// healthCheck handles health check requests
func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"version": "v0.0.1",
		"service": "probe-agent",
	})
}

// handlePolicyEvent processes policy events from the operator
func handlePolicyEvent(c *gin.Context) {
	var event PolicyEventRequest

	// Bind and validate JSON request
	if err := c.ShouldBindJSON(&event); err != nil {
		log.Printf("Invalid request: %v", err)
		c.JSON(http.StatusBadRequest, EventResponse{
			Status:  "error",
			Message: fmt.Sprintf("Invalid request: %v", err),
		})
		return
	}

	// Log the event details
	log.Println("==================== POLICY EVENT RECEIVED ====================")
	log.Printf("Action:        %s", event.Action)
	log.Printf("Resource Name: %s", event.ResourceName)
	log.Printf("Namespace:     %s", event.Namespace)
	log.Println("Policy Details:")
	log.Printf("  - Name:        %s", event.Policy.Name)
	log.Printf("  - Type:        %s", event.Policy.Type)
	log.Printf("  - Target:      %s", event.Policy.Target)
	log.Printf("  - Enabled:     %v", event.Policy.Enabled)
	if event.Policy.Description != "" {
		log.Printf("  - Description: %s", event.Policy.Description)
	}
	if len(event.Policy.Parameters) > 0 {
		log.Println("  - Parameters:")
		for key, value := range event.Policy.Parameters {
			log.Printf("      %s: %s", key, value)
		}
	}
	log.Println("===============================================================")

	// TODO: Execute the actual policy based on the event
	// For now, just acknowledge receipt
	c.JSON(http.StatusOK, EventResponse{
		Status:  "success",
		Message: fmt.Sprintf("Policy event '%s' received for resource '%s'", event.Action, event.ResourceName),
	})
}
