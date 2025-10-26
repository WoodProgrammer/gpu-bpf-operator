/*
Copyright 2025 WoodProgrammer.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	gpuv1alpha1 "github.com/WoodProgrammer/gpu-bpf-operator/api/v1alpha1"
)

const (
	finalizerName      = "gpu.obs.gpu/finalizer"
	monitoringEndpoint = "http://dummy.monitoring.gpu.svc:9090/reconfig"
	configFilePath     = "CONFIG.md"
)

// CudaEBPFPolicyReconciler reconciles a CudaEBPFPolicy object
type CudaEBPFPolicyReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// PolicyConfig represents the configuration from CONFIG.md
type PolicyConfig struct {
	Hash     string         `json:"hash"`
	Policies []PolicyDetail `json:"policies"`
}

// PolicyDetail represents a single policy in the configuration
type PolicyDetail struct {
	ID           string                 `json:"id"`
	LibPath      string                 `json:"libPath"`
	Mode         string                 `json:"mode"`
	ProcessRegex string                 `json:"processRegex"`
	Functions    []gpuv1alpha1.Function `json:"functions"`
	Output       map[string]interface{} `json:"output"`
}

// ReconfigRequest represents the request sent to the monitoring server
type ReconfigRequest struct {
	Action string       `json:"action"`
	Policy PolicyDetail `json:"policy"`
}

// +kubebuilder:rbac:groups=gpu.obs.gpu,resources=cudaebpfpolicies,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=gpu.obs.gpu,resources=cudaebpfpolicies/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=gpu.obs.gpu,resources=cudaebpfpolicies/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *CudaEBPFPolicyReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	// Fetch the CudaEBPFPolicy instance
	policy := &gpuv1alpha1.CudaEBPFPolicy{}
	err := r.Get(ctx, req.NamespacedName, policy)
	if err != nil {
		if errors.IsNotFound(err) {
			// Object not found, could have been deleted after reconcile request
			log.Info("CudaEBPFPolicy resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
		log.Error(err, "Failed to get CudaEBPFPolicy")
		return ctrl.Result{}, err
	}

	// Determine the action based on the resource state
	action := ""

	// Check if the resource is being deleted
	if !policy.DeletionTimestamp.IsZero() {
		action = "delete"
		if controllerutil.ContainsFinalizer(policy, finalizerName) {
			// Send delete action to monitoring server
			if err := r.sendReconfigRequest(ctx, action, policy); err != nil {
				log.Error(err, "Failed to send delete request to monitoring server")
				return ctrl.Result{}, err
			}

			// Remove finalizer
			controllerutil.RemoveFinalizer(policy, finalizerName)
			if err := r.Update(ctx, policy); err != nil {
				return ctrl.Result{}, err
			}
		}
		return ctrl.Result{}, nil
	}

	// Add finalizer if it doesn't exist
	if !controllerutil.ContainsFinalizer(policy, finalizerName) {
		controllerutil.AddFinalizer(policy, finalizerName)
		if err := r.Update(ctx, policy); err != nil {
			return ctrl.Result{}, err
		}
	}

	// Calculate hash of current spec
	currentHash, err := r.calculateHash(policy)
	if err != nil {
		log.Error(err, "Failed to calculate hash")
		return ctrl.Result{}, err
	}

	// Determine if this is an add or update
	if policy.Status.ObservedHash == "" {
		action = "add"
	} else if policy.Status.ObservedHash != currentHash {
		action = "update"
	} else {
		// No changes detected
		log.Info("No changes detected in policy spec")
		return ctrl.Result{}, nil
	}

	// Send reconfig request to monitoring server
	if err := r.sendReconfigRequest(ctx, action, policy); err != nil {
		log.Error(err, "Failed to send reconfig request to monitoring server")
		return ctrl.Result{RequeueAfter: 30 * time.Second}, err
	}

	// Update status with new hash
	policy.Status.ObservedHash = currentHash
	if err := r.Status().Update(ctx, policy); err != nil {
		log.Error(err, "Failed to update policy status")
		return ctrl.Result{}, err
	}

	log.Info("Successfully reconciled CudaEBPFPolicy", "action", action, "policy", req.NamespacedName)
	return ctrl.Result{}, nil
}

// calculateHash computes a hash of the policy spec to detect changes
func (r *CudaEBPFPolicyReconciler) calculateHash(policy *gpuv1alpha1.CudaEBPFPolicy) (string, error) {
	specBytes, err := json.Marshal(policy.Spec)
	if err != nil {
		return "", err
	}
	hash := sha256.Sum256(specBytes)
	return fmt.Sprintf("%x", hash), nil
}

// sendReconfigRequest sends a reconfiguration request to the monitoring server
func (r *CudaEBPFPolicyReconciler) sendReconfigRequest(ctx context.Context, action string, policy *gpuv1alpha1.CudaEBPFPolicy) error {
	log := logf.FromContext(ctx)

	// Build policy detail from CRD and CONFIG.md
	policyDetail, err := r.buildPolicyDetail(policy)
	if err != nil {
		return fmt.Errorf("failed to build policy detail: %w", err)
	}

	// Create reconfig request
	req := ReconfigRequest{
		Action: action,
		Policy: policyDetail,
	}

	// Marshal request to JSON
	reqBody, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal reconfig request: %w", err)
	}

	// Send HTTP POST request
	httpReq, err := http.NewRequestWithContext(ctx, "POST", monitoringEndpoint, bytes.NewBuffer(reqBody))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(httpReq)
	if err != nil {
		return fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("monitoring server returned error status %d: %s", resp.StatusCode, string(body))
	}

	log.Info("Successfully sent reconfig request", "action", action, "endpoint", monitoringEndpoint)
	return nil
}

// buildPolicyDetail constructs a PolicyDetail from the CRD and CONFIG.md
func (r *CudaEBPFPolicyReconciler) buildPolicyDetail(policy *gpuv1alpha1.CudaEBPFPolicy) (PolicyDetail, error) {
	// Read CONFIG.md for additional configuration
	config, err := r.readConfig()
	if err != nil {
		// If config file doesn't exist, use CRD spec only
		return PolicyDetail{
			ID:           fmt.Sprintf("%s@%s", policy.Name, policy.Namespace),
			LibPath:      policy.Spec.LibPath,
			Mode:         policy.Spec.Mode,
			ProcessRegex: policy.Spec.ProcessRegex,
			Functions:    policy.Spec.Functions,
			Output: map[string]interface{}{
				"format": policy.Spec.OutputFormat,
			},
		}, nil
	}

	// Merge CONFIG.md with CRD spec (CRD takes precedence)
	policyDetail := PolicyDetail{
		ID:           fmt.Sprintf("%s@%s", policy.Name, policy.Namespace),
		LibPath:      policy.Spec.LibPath,
		Mode:         policy.Spec.Mode,
		ProcessRegex: policy.Spec.ProcessRegex,
		Functions:    policy.Spec.Functions,
	}

	// Use output format from CONFIG.md if not specified in CRD
	if len(config.Policies) > 0 && policy.Spec.OutputFormat == "" {
		policyDetail.Output = config.Policies[0].Output
	} else {
		policyDetail.Output = map[string]interface{}{
			"format": policy.Spec.OutputFormat,
		}
	}

	return policyDetail, nil
}

// readConfig reads and parses the CONFIG.md file
func (r *CudaEBPFPolicyReconciler) readConfig() (*PolicyConfig, error) {
	data, err := os.ReadFile(configFilePath)
	if err != nil {
		return nil, err
	}

	var config PolicyConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse CONFIG.md: %w", err)
	}

	return &config, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *CudaEBPFPolicyReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&gpuv1alpha1.CudaEBPFPolicy{}).
		Named("cudaebpfpolicy").
		Complete(r)
}
