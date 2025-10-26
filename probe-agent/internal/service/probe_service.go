package service

import (
	"fmt"
	"gpu-bpf-operator/probe-agent/internal/model"
	"log"
	"sync"
)

// ProbeService defines the interface for probe execution business logic
type ProbeService interface {
	CreateProbeExecution(req *model.ProbeExecutionRequest) error
	DeleteProbeExecution(resourceName, namespace string) error
	GetProbeExecutionStatus(resourceName, namespace string) (*model.ProbeExecutionStatus, error)
}

// probeService implements ProbeService interface
type probeService struct {
	mu         sync.RWMutex
	executions map[string]*model.ProbeExecutionStatus // key: namespace/resourceName
}

// NewProbeService creates a new instance of ProbeService
func NewProbeService() ProbeService {
	return &probeService{
		executions: make(map[string]*model.ProbeExecutionStatus),
	}
}

// CreateProbeExecution handles the creation of a new probe execution
func (s *probeService) CreateProbeExecution(req *model.ProbeExecutionRequest) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := fmt.Sprintf("%s/%s", req.Namespace, req.ResourceName)

	// Log the policy details
	log.Println("==================== CREATE PROBE EXECUTION ====================")
	log.Printf("Action:        %s", req.Action)
	log.Printf("Resource Name: %s", req.ResourceName)
	log.Printf("Namespace:     %s", req.Namespace)
	log.Println("Policy Details:")
	log.Printf("  - Name:        %s", req.Policy.Name)
	log.Printf("  - Type:        %s", req.Policy.Type)
	log.Printf("  - Target:      %s", req.Policy.Target)
	log.Printf("  - Enabled:     %v", req.Policy.Enabled)
	if req.Policy.Description != "" {
		log.Printf("  - Description: %s", req.Policy.Description)
	}
	if len(req.Policy.Parameters) > 0 {
		log.Println("  - Parameters:")
		for key, value := range req.Policy.Parameters {
			log.Printf("      %s: %s", key, value)
		}
	}
	log.Println("===============================================================")

	// TODO: Implement actual eBPF probe loading logic here
	// For now, just store the execution status
	s.executions[key] = &model.ProbeExecutionStatus{
		ResourceName: req.ResourceName,
		Namespace:    req.Namespace,
		Status:       "running",
		Message:      fmt.Sprintf("Probe execution created for policy '%s'", req.Policy.Name),
		Policy:       req.Policy.Name,
	}

	return nil
}

// DeleteProbeExecution handles the deletion of a probe execution
func (s *probeService) DeleteProbeExecution(resourceName, namespace string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := fmt.Sprintf("%s/%s", namespace, resourceName)

	if _, exists := s.executions[key]; !exists {
		return fmt.Errorf("probe execution not found for resource '%s' in namespace '%s'", resourceName, namespace)
	}

	log.Printf("==================== DELETE PROBE EXECUTION ====================")
	log.Printf("Resource Name: %s", resourceName)
	log.Printf("Namespace:     %s", namespace)
	log.Println("===============================================================")

	// TODO: Implement actual eBPF probe unloading logic here
	//delete(s.executions, key)

	return nil
}

// GetProbeExecutionStatus retrieves the status of a probe execution
func (s *probeService) GetProbeExecutionStatus(resourceName, namespace string) (*model.ProbeExecutionStatus, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	key := fmt.Sprintf("%s/%s", namespace, resourceName)

	status, exists := s.executions[key]
	if !exists {
		return nil, fmt.Errorf("probe execution not found for resource '%s' in namespace '%s'", resourceName, namespace)
	}

	return status, nil
}
