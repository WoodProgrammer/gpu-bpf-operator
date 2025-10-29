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
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

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
			// Daemonset logic
			// Remove finalizer
			found := &appsv1.DaemonSet{}
			err = r.Get(ctx, types.NamespacedName{Name: policy.Name, Namespace: policy.Namespace}, found)
			if err != nil && errors.IsNotFound(err) {
				log.Error(err, "Failed to create new Daemonset", "Daemonset.Namespace", found.Namespace, "Daemonset.Name", found.Name)
				r.Delete(ctx, found)
				return ctrl.Result{}, err
			}

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
		log.Info("Update detected on policy definitions")

		updatedHash, err := r.calculateHash(policy)
		policy.Status.ObservedHash = updatedHash
		log.Info("Hash values adjusted")
		if err != nil {
			log.Error(err, "Failed to calculate hash")
			return ctrl.Result{}, err
		}

		ds := r.createDaemonsetProbeAgent(policy)
		log.Info("Update a new Daemonset", "Daemonset.Namespace", ds.Namespace, "Daemonset.Name", ds.Name)
		err = r.Update(ctx, ds)
		if err != nil {
			log.Error(err, "Failed to update new Daemonset", "Daemonset.Namespace", ds.Namespace, "Daemonset.Name", ds.Name)
			return ctrl.Result{}, err
		}
		// Daemonset created successfully - return and requeue
		return ctrl.Result{Requeue: true}, nil

	} else {
		// No changes detected
		log.Info("No changes detected in policy spec")
		return ctrl.Result{}, nil
	}
	// Check if the deployment already exists, if not create a new one
	found := &appsv1.DaemonSet{}
	err = r.Get(ctx, types.NamespacedName{Name: policy.Name, Namespace: policy.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		// Define a new deployment
		ds := r.createDaemonsetProbeAgent(policy)
		log.Info("Creating a new Daemonset", "Daemonset.Namespace", ds.Namespace, "Daemonset.Name", ds.Name)
		err = r.Create(ctx, ds)
		if err != nil {
			log.Error(err, "Failed to create new Daemonset", "Daemonset.Namespace", ds.Namespace, "Daemonset.Name", ds.Name)
			return ctrl.Result{}, err
		}
		// Deployment created successfully - return and requeue
		return ctrl.Result{Requeue: true}, nil
	} else if err != nil {
		log.Error(err, "Failed to get Daemonset")
		return ctrl.Result{}, err
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

func (r *CudaEBPFPolicyReconciler) createDaemonsetProbeAgent(policy *gpuv1alpha1.CudaEBPFPolicy) *appsv1.DaemonSet {
	labels := map[string]string{
		"app": "gpu-operator",
	}
	ProbeCallsDetails := r.EncodeProbeCalls(policy)
	fmt.Println("The probe call details are", ProbeCallsDetails)
	ds := &appsv1.DaemonSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      policy.Name,
			Namespace: "gpu-bpf-operator", // TODO namespace should be optional
		},
		Spec: appsv1.DaemonSetSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Image: policy.Spec.Image,
						Name:  "bpf-tracer-agent",
						Ports: []corev1.ContainerPort{{
							ContainerPort: 9090,
							Name:          "bpfpolicyagent",
						}},
						Env: []corev1.EnvVar{{
							Name:  "LIB_PATH",
							Value: policy.Spec.LibPath,
						},
							{
								Name:  "PROBE_CALLS",
								Value: ProbeCallsDetails,
							}},
					}},
				},
			},
		},
	}
	// Set Memcached instance as the owner and controller
	ctrl.SetControllerReference(policy, ds, r.Scheme)
	return ds
}

func (r *CudaEBPFPolicyReconciler) EncodeProbeCalls(policy *gpuv1alpha1.CudaEBPFPolicy) string {

	return ""
}

// SetupWithManager sets up the controller with the Manager.
func (r *CudaEBPFPolicyReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&gpuv1alpha1.CudaEBPFPolicy{}).
		Named("cudaebpfpolicy").
		Complete(r)
}
