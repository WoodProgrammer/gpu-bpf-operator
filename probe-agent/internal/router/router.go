package router

import (
	"gpu-bpf-operator/probe-agent/internal/handler"
	"gpu-bpf-operator/probe-agent/internal/service"

	"github.com/gin-gonic/gin"
)

const (
	// API version prefix
	apiV1 = "/v1"

	// Endpoints
	healthEndpoint = "/health"
)

// Setup initializes and configures all routes
func Setup(version string) *gin.Engine {
	router := gin.Default()

	// Initialize services
	probeService := service.NewProbeService()

	// Initialize handlers
	probeHandler := handler.NewProbeHandler(probeService)
	healthHandler := handler.NewHealthHandler(version)

	// Health check endpoint
	router.GET(healthEndpoint, healthHandler.HealthCheck)

	// API v1 routes
	v1 := router.Group(apiV1)
	{
		// Probe execution management
		v1.POST("/reconfig", probeHandler.CreateProbeExecution)
		v1.DELETE("/reconfig", probeHandler.DeleteProbeExecution)
		v1.GET("/reconfig", probeHandler.GetProbeExecutionStatus)
	}

	return router
}
