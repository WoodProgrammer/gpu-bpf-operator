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

package v1alpha1

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	gpuv1alpha1 "github.com/WoodProgrammer/gpu-bpf-operator/api/v1alpha1"
)

// nolint:unused
// log is for logging in this package.
var cudaebpfpolicylog = logf.Log.WithName("cudaebpfpolicy-resource")

// SetupCudaEBPFPolicyWebhookWithManager registers the webhook for CudaEBPFPolicy in the manager.
func SetupCudaEBPFPolicyWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).For(&gpuv1alpha1.CudaEBPFPolicy{}).
		WithValidator(&CudaEBPFPolicyCustomValidator{}).
		WithDefaulter(&CudaEBPFPolicyCustomDefaulter{}).
		Complete()
}

// TODO(user): EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

// +kubebuilder:webhook:path=/mutate-gpu-obs-gpu-v1alpha1-cudaebpfpolicy,mutating=true,failurePolicy=fail,sideEffects=None,groups=gpu.obs.gpu,resources=cudaebpfpolicies,verbs=create;update,versions=v1alpha1,name=mcudaebpfpolicy-v1alpha1.kb.io,admissionReviewVersions=v1

// CudaEBPFPolicyCustomDefaulter struct is responsible for setting default values on the custom resource of the
// Kind CudaEBPFPolicy when those are created or updated.
//
// NOTE: The +kubebuilder:object:generate=false marker prevents controller-gen from generating DeepCopy methods,
// as it is used only for temporary operations and does not need to be deeply copied.
type CudaEBPFPolicyCustomDefaulter struct {
	// TODO(user): Add more fields as needed for defaulting
}

var _ webhook.CustomDefaulter = &CudaEBPFPolicyCustomDefaulter{}

// Default implements webhook.CustomDefaulter so a webhook will be registered for the Kind CudaEBPFPolicy.
func (d *CudaEBPFPolicyCustomDefaulter) Default(_ context.Context, obj runtime.Object) error {
	cudaebpfpolicy, ok := obj.(*gpuv1alpha1.CudaEBPFPolicy)

	if !ok {
		return fmt.Errorf("expected an CudaEBPFPolicy object but got %T", obj)
	}
	cudaebpfpolicylog.Info("Defaulting for CudaEBPFPolicy", "name", cudaebpfpolicy.GetName())

	// TODO(user): fill in your defaulting logic.

	return nil
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
// NOTE: The 'path' attribute must follow a specific pattern and should not be modified directly here.
// Modifying the path for an invalid path can cause API server errors; failing to locate the webhook.
// +kubebuilder:webhook:path=/validate-gpu-obs-gpu-v1alpha1-cudaebpfpolicy,mutating=false,failurePolicy=fail,sideEffects=None,groups=gpu.obs.gpu,resources=cudaebpfpolicies,verbs=create;update,versions=v1alpha1,name=vcudaebpfpolicy-v1alpha1.kb.io,admissionReviewVersions=v1

// CudaEBPFPolicyCustomValidator struct is responsible for validating the CudaEBPFPolicy resource
// when it is created, updated, or deleted.
//
// NOTE: The +kubebuilder:object:generate=false marker prevents controller-gen from generating DeepCopy methods,
// as this struct is used only for temporary operations and does not need to be deeply copied.
type CudaEBPFPolicyCustomValidator struct {
	// TODO(user): Add more fields as needed for validation
}

var _ webhook.CustomValidator = &CudaEBPFPolicyCustomValidator{}

// ValidateCreate implements webhook.CustomValidator so a webhook will be registered for the type CudaEBPFPolicy.
func (v *CudaEBPFPolicyCustomValidator) ValidateCreate(_ context.Context, obj runtime.Object) (admission.Warnings, error) {
	cudaebpfpolicy, ok := obj.(*gpuv1alpha1.CudaEBPFPolicy)
	if !ok {
		return nil, fmt.Errorf("expected a CudaEBPFPolicy object but got %T", obj)
	}
	cudaebpfpolicylog.Info("Validation for CudaEBPFPolicy upon creation", "name", cudaebpfpolicy.GetName())

	// TODO(user): fill in your validation logic upon object creation.

	return nil, nil
}

// ValidateUpdate implements webhook.CustomValidator so a webhook will be registered for the type CudaEBPFPolicy.
func (v *CudaEBPFPolicyCustomValidator) ValidateUpdate(_ context.Context, oldObj, newObj runtime.Object) (admission.Warnings, error) {
	cudaebpfpolicy, ok := newObj.(*gpuv1alpha1.CudaEBPFPolicy)
	if !ok {
		return nil, fmt.Errorf("expected a CudaEBPFPolicy object for the newObj but got %T", newObj)
	}
	cudaebpfpolicylog.Info("Validation for CudaEBPFPolicy upon update", "name", cudaebpfpolicy.GetName())

	// TODO(user): fill in your validation logic upon object update.

	return nil, nil
}

// ValidateDelete implements webhook.CustomValidator so a webhook will be registered for the type CudaEBPFPolicy.
func (v *CudaEBPFPolicyCustomValidator) ValidateDelete(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	cudaebpfpolicy, ok := obj.(*gpuv1alpha1.CudaEBPFPolicy)
	if !ok {
		return nil, fmt.Errorf("expected a CudaEBPFPolicy object but got %T", obj)
	}
	cudaebpfpolicylog.Info("Validation for CudaEBPFPolicy upon deletion", "name", cudaebpfpolicy.GetName())

	// TODO(user): fill in your validation logic upon object deletion.

	return nil, nil
}
