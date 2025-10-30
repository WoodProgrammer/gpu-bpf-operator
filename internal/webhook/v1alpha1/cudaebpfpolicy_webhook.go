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
	"strings"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
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

	// Set default mode to pidwatch if not specified
	if cudaebpfpolicy.Spec.Mode == "" {
		cudaebpfpolicy.Spec.Mode = "pidwatch"
	}

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

	return nil, v.validateCudaEBPFPolicy(cudaebpfpolicy)
}

// ValidateUpdate implements webhook.CustomValidator so a webhook will be registered for the type CudaEBPFPolicy.
func (v *CudaEBPFPolicyCustomValidator) ValidateUpdate(_ context.Context, oldObj, newObj runtime.Object) (admission.Warnings, error) {
	cudaebpfpolicy, ok := newObj.(*gpuv1alpha1.CudaEBPFPolicy)
	if !ok {
		return nil, fmt.Errorf("expected a CudaEBPFPolicy object for the newObj but got %T", newObj)
	}
	cudaebpfpolicylog.Info("Validation for CudaEBPFPolicy upon update", "name", cudaebpfpolicy.GetName())

	return nil, v.validateCudaEBPFPolicy(cudaebpfpolicy)
}

// ValidateDelete implements webhook.CustomValidator so a webhook will be registered for the type CudaEBPFPolicy.
func (v *CudaEBPFPolicyCustomValidator) ValidateDelete(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	cudaebpfpolicy, ok := obj.(*gpuv1alpha1.CudaEBPFPolicy)
	if !ok {
		return nil, fmt.Errorf("expected a CudaEBPFPolicy object but got %T", obj)
	}
	cudaebpfpolicylog.Info("Validation for CudaEBPFPolicy upon deletion", "name", cudaebpfpolicy.GetName())

	// No validation needed on delete
	return nil, nil
}

// validateCudaEBPFPolicy validates the CudaEBPFPolicy spec
func (v *CudaEBPFPolicyCustomValidator) validateCudaEBPFPolicy(policy *gpuv1alpha1.CudaEBPFPolicy) error {
	var allErrs field.ErrorList

	// Validate functions field
	if err := v.validateFunctions(policy.Spec.Functions, field.NewPath("spec").Child("functions")); err != nil {
		allErrs = append(allErrs, err...)
	}

	// Validate mode field
	if err := v.validateMode(policy.Spec.Mode, field.NewPath("spec").Child("mode")); err != nil {
		allErrs = append(allErrs, err)
	}

	// Validate outputFormat field
	if err := v.validateOutputFormat(policy.Spec.OutputFormat, field.NewPath("spec").Child("output")); err != nil {
		allErrs = append(allErrs, err)
	}

	// Validate libPath is not empty
	if policy.Spec.LibPath == "" {
		allErrs = append(allErrs, field.Required(field.NewPath("spec").Child("libPath"), "libPath must be specified"))
	}

	// Validate image is not empty
	if policy.Spec.Image == "" {
		allErrs = append(allErrs, field.Required(field.NewPath("spec").Child("image"), "image must be specified"))
	}

	if len(allErrs) == 0 {
		return nil
	}

	return allErrs.ToAggregate()
}

// validateFunctions validates the functions field
func (v *CudaEBPFPolicyCustomValidator) validateFunctions(functions []gpuv1alpha1.Function, fldPath *field.Path) field.ErrorList {
	var allErrs field.ErrorList

	// Validate functions array is not empty
	if len(functions) == 0 {
		allErrs = append(allErrs, field.Required(fldPath, "at least one function must be specified"))
		return allErrs
	}

	// Track function names to detect duplicates
	functionNames := make(map[string]bool)

	for i, fn := range functions {
		funcPath := fldPath.Index(i)

		// Validate function name is not empty
		if fn.Name == "" {
			allErrs = append(allErrs, field.Required(funcPath.Child("name"), "function name must be specified"))
		}

		// Check for duplicate function names
		if functionNames[fn.Name] {
			allErrs = append(allErrs, field.Duplicate(funcPath.Child("name"), fn.Name))
		}
		functionNames[fn.Name] = true

		// Validate function kind
		validKinds := []string{"uprobe", "uretprobe", "kprobe", "kretprobe"}
		if !contains(validKinds, fn.Kind) {
			allErrs = append(allErrs, field.NotSupported(funcPath.Child("kind"), fn.Kind, validKinds))
		}

		// Validate arguments if present
		if len(fn.Args) > 0 {
			if err := v.validateArgs(fn.Args, funcPath.Child("args")); err != nil {
				allErrs = append(allErrs, err...)
			}
		}
	}

	return allErrs
}

// validateArgs validates function arguments
func (v *CudaEBPFPolicyCustomValidator) validateArgs(args []gpuv1alpha1.Arg, fldPath *field.Path) field.ErrorList {
	var allErrs field.ErrorList

	// Track argument indices to detect duplicates
	argIndices := make(map[int]bool)

	for i, arg := range args {
		argPath := fldPath.Index(i)

		// Validate argument name is not empty
		if arg.Name == "" {
			allErrs = append(allErrs, field.Required(argPath.Child("name"), "argument name must be specified"))
		}

		// Validate argument index is non-negative
		if arg.Index < 0 {
			allErrs = append(allErrs, field.Invalid(argPath.Child("index"), arg.Index, "argument index must be non-negative"))
		}

		// Check for duplicate argument indices
		if argIndices[arg.Index] {
			allErrs = append(allErrs, field.Duplicate(argPath.Child("index"), arg.Index))
		}
		argIndices[arg.Index] = true
	}

	return allErrs
}

// validateMode validates the mode field
func (v *CudaEBPFPolicyCustomValidator) validateMode(mode string, fldPath *field.Path) *field.Error {
	validModes := []string{"pidwatch", "systemwide"}
	if !contains(validModes, mode) {
		return field.NotSupported(fldPath, mode, validModes)
	}
	return nil
}

// validateOutputFormat validates the output format field
func (v *CudaEBPFPolicyCustomValidator) validateOutputFormat(format string, fldPath *field.Path) *field.Error {
	// Empty format is allowed (will be defaulted)
	if format == "" {
		return nil
	}

	validFormats := []string{"ndjson", "prometheus"}
	if !contains(validFormats, format) {
		return field.NotSupported(fldPath, format, validFormats)
	}
	return nil
}

// contains checks if a string is in a slice
func contains(slice []string, str string) bool {
	for _, s := range slice {
		if strings.EqualFold(s, str) {
			return true
		}
	}
	return false
}
