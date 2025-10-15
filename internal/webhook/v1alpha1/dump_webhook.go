/*
Copyright 2025.

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
	"errors"
	"fmt"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	dumpv1alpha1 "github.com/WoodProgrammer/kubexdp-operator/api/v1alpha1"
)

// nolint:unused
// log is for logging in this package.
var dumplog = logf.Log.WithName("dump-resource")

// SetupDumpWebhookWithManager registers the webhook for Dump in the manager.
func SetupDumpWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).For(&dumpv1alpha1.Dump{}).
		WithValidator(&DumpCustomValidator{}).
		WithDefaulter(&DumpCustomDefaulter{}).
		Complete()
}

// TODO(user): EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

// +kubebuilder:webhook:path=/mutate-dump-kubexdp-io-v1alpha1-dump,mutating=true,failurePolicy=fail,sideEffects=None,groups=dump.kubexdp.io,resources=dumps,verbs=create;update,versions=v1alpha1,name=mdump-v1alpha1.kb.io,admissionReviewVersions=v1

// DumpCustomDefaulter struct is responsible for setting default values on the custom resource of the
// Kind Dump when those are created or updated.
//
// NOTE: The +kubebuilder:object:generate=false marker prevents controller-gen from generating DeepCopy methods,
// as it is used only for temporary operations and does not need to be deeply copied.
type DumpCustomDefaulter struct {
	// TODO(user): Add more fields as needed for defaulting
}

var _ webhook.CustomDefaulter = &DumpCustomDefaulter{}

// Default implements webhook.CustomDefaulter so a webhook will be registered for the Kind Dump.
func (d *DumpCustomDefaulter) Default(_ context.Context, obj runtime.Object) error {
	dump, ok := obj.(*dumpv1alpha1.Dump)

	if !ok {
		return fmt.Errorf("expected an Dump object but got %T", obj)
	}
	dumplog.Info("Defaulting for Dump", "name", dump.GetName())

	// TODO(user): fill in your defaulting logic.

	return nil
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
// NOTE: The 'path' attribute must follow a specific pattern and should not be modified directly here.
// Modifying the path for an invalid path can cause API server errors; failing to locate the webhook.
// +kubebuilder:webhook:path=/validate-dump-kubexdp-io-v1alpha1-dump,mutating=false,failurePolicy=fail,sideEffects=None,groups=dump.kubexdp.io,resources=dumps,verbs=create;update,versions=v1alpha1,name=vdump-v1alpha1.kb.io,admissionReviewVersions=v1

// DumpCustomValidator struct is responsible for validating the Dump resource
// when it is created, updated, or deleted.
//
// NOTE: The +kubebuilder:object:generate=false marker prevents controller-gen from generating DeepCopy methods,
// as this struct is used only for temporary operations and does not need to be deeply copied.
type DumpCustomValidator struct {
	// TODO(user): Add more fields as needed for validation
}

var _ webhook.CustomValidator = &DumpCustomValidator{}

// ValidateCreate implements webhook.CustomValidator so a webhook will be registered for the type Dump.
func (v *DumpCustomValidator) ValidateCreate(_ context.Context, obj runtime.Object) (admission.Warnings, error) {
	var allErrs field.ErrorList

	dump, ok := obj.(*dumpv1alpha1.Dump)
	if !ok {
		return nil, fmt.Errorf("expected a Dump object but got %T", obj)
	}

	if dump.Spec.Schedule == "* * * * *" {
		err := errors.New("the frequency error on schedule definition of dump objects")
		allErrs = append(allErrs, field.Invalid(nil, dump.Spec.Schedule, err.Error()))

		dumplog.Info("Validation error", "name", dump.GetName())

		return nil, apierrors.NewInvalid(
			schema.GroupKind{Group: "dump.kubexdp.io/v1alpha1", Kind: "Dump"},
			dump.Name, allErrs)
	}
	dumplog.Info("Validation for Dump upon creation", "name", dump.GetName())

	// TODO(user): fill in your validation logic upon object creation.

	return nil, nil
}

// ValidateUpdate implements webhook.CustomValidator so a webhook will be registered for the type Dump.
func (v *DumpCustomValidator) ValidateUpdate(_ context.Context, oldObj, newObj runtime.Object) (admission.Warnings, error) {
	dump, ok := newObj.(*dumpv1alpha1.Dump)
	if !ok {
		return nil, fmt.Errorf("expected a Dump object for the newObj but got %T", newObj)
	}
	dumplog.Info("Validation for Dump upon update", "name", dump.GetName())

	// TODO(user): fill in your validation logic upon object update.

	return nil, nil
}

// ValidateDelete implements webhook.CustomValidator so a webhook will be registered for the type Dump.
func (v *DumpCustomValidator) ValidateDelete(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	dump, ok := obj.(*dumpv1alpha1.Dump)
	if !ok {
		return nil, fmt.Errorf("expected a Dump object but got %T", obj)
	}
	dumplog.Info("Validation for Dump upon deletion", "name", dump.GetName())

	// TODO(user): fill in your validation logic upon object deletion.

	return nil, nil
}
