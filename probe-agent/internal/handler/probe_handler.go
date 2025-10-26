package handler

import (
	"gpu-bpf-operator/probe-agent/internal/model"
	"gpu-bpf-operator/probe-agent/internal/service"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ProbeHandler handles HTTP requests for probe execution management
type ProbeHandler struct {
	probeService service.ProbeService
}

// NewProbeHandler creates a new ProbeHandler instance
func NewProbeHandler(probeService service.ProbeService) *ProbeHandler {
	return &ProbeHandler{
		probeService: probeService,
	}
}

// CreateProbeExecution handles POST requests to create a new probe execution
func (h *ProbeHandler) CreateProbeExecution(c *gin.Context) {
	var req model.ProbeExecutionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Invalid request: %v", err)
		c.JSON(http.StatusBadRequest, model.APIResponse{
			Status:  "error",
			Message: "Invalid request body: " + err.Error(),
		})
		return
	}

	if err := h.probeService.CreateProbeExecution(&req); err != nil {
		log.Printf("Failed to create probe execution: %v", err)
		c.JSON(http.StatusInternalServerError, model.APIResponse{
			Status:  "error",
			Message: "Failed to create probe execution: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, model.APIResponse{
		Status:  "success",
		Message: "Probe execution created successfully",
	})
}

// DeleteProbeExecution handles DELETE requests to remove a probe execution
func (h *ProbeHandler) DeleteProbeExecution(c *gin.Context) {
	resourceName := c.Query("resource_name")
	namespace := c.Query("namespace")

	if resourceName == "" || namespace == "" {
		c.JSON(http.StatusBadRequest, model.APIResponse{
			Status:  "error",
			Message: "Missing required query parameters: resource_name and namespace",
		})
		return
	}

	if err := h.probeService.DeleteProbeExecution(resourceName, namespace); err != nil {
		log.Printf("Failed to delete probe execution: %v", err)
		c.JSON(http.StatusNotFound, model.APIResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.APIResponse{
		Status:  "success",
		Message: "Probe execution deleted successfully",
	})
}

// GetProbeExecutionStatus handles GET requests to retrieve probe execution status
func (h *ProbeHandler) GetProbeExecutionStatus(c *gin.Context) {
	resourceName := c.Query("resource_name")
	namespace := c.Query("namespace")

	if resourceName == "" || namespace == "" {
		c.JSON(http.StatusBadRequest, model.APIResponse{
			Status:  "error",
			Message: "Missing required query parameters: resource_name and namespace",
		})
		return
	}

	status, err := h.probeService.GetProbeExecutionStatus(resourceName, namespace)
	if err != nil {
		log.Printf("Failed to get probe execution status: %v", err)
		c.JSON(http.StatusNotFound, model.APIResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.APIResponse{
		Status:  "success",
		Message: "Probe execution status retrieved successfully",
		Data:    status,
	})
}
