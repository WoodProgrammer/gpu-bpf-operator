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

package controller

import (
	"context"

	dumpv1alpha1 "github.com/WoodProgrammer/kubexdp-operator/api/v1alpha1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

// DumpReconciler reconciles a Dump object
type DumpReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=dump.kubexdp.io,resources=dumps,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=dump.kubexdp.io,resources=dumps/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=dump.kubexdp.io,resources=dumps/finalizers,verbs=update
// +kubebuilder:rbac:groups=batch,resources=cronjobs,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Dump object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.21.0/pkg/reconcile
func (r *DumpReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	// Fetch the Dump instance
	dump := &dumpv1alpha1.Dump{}
	err := r.Get(ctx, req.NamespacedName, dump)
	if err != nil {
		if client.IgnoreNotFound(err) == nil {
			// Object not found, it was deleted
			log.Info("DELETE event detected", "name", req.Name, "namespace", req.Namespace)
			// Delete associated CronJob
			cronJob := &batchv1.CronJob{}
			cronJobName := req.Name + "-cronjob"
			cronJobErr := r.Get(ctx, client.ObjectKey{Name: cronJobName, Namespace: req.Namespace}, cronJob)
			if cronJobErr == nil {
				// CronJob exists, delete it
				if err := r.Delete(ctx, cronJob); err != nil {
					log.Error(err, "Failed to delete CronJob", "cronjob", cronJobName)
					return ctrl.Result{}, err
				}
				log.Info("Successfully deleted CronJob", "cronjob", cronJobName)
			}
			return ctrl.Result{}, nil
		}
		// Error reading the object
		log.Error(err, "Failed to get Dump resource")
		return ctrl.Result{}, err
	}

	// Check if object is being deleted (has deletion timestamp)
	if !dump.ObjectMeta.DeletionTimestamp.IsZero() {
		log.Info("DELETE event detected (finalizer)", "name", dump.Name, "namespace", dump.Namespace)
		// Delete associated CronJob
		cronJob := &batchv1.CronJob{}
		cronJobName := dump.Name + "-cronjob"
		cronJobErr := r.Get(ctx, client.ObjectKey{Name: cronJobName, Namespace: dump.Namespace}, cronJob)
		if cronJobErr == nil {
			// CronJob exists, delete it
			if err := r.Delete(ctx, cronJob); err != nil {
				log.Error(err, "Failed to delete CronJob", "cronjob", cronJobName)
				return ctrl.Result{}, err
			}
			log.Info("Successfully deleted CronJob", "cronjob", cronJobName)
		}
		return ctrl.Result{}, nil
	}

	// Determine if this is CREATE or UPDATE
	if dump.Status.ObservedGeneration == 0 {
		// This is a CREATE event (first reconciliation)
		log.Info("CREATE event detected", "name", dump.Name, "namespace", dump.Namespace)
	} else {
		// This is an UPDATE event
		log.Info("UPDATE event detected", "name", dump.Name, "namespace", dump.Namespace, "generation", dump.Generation)
	}

	// Create or Update CronJob
	cronJob := &batchv1.CronJob{}
	cronJobName := dump.Name + "-cronjob"
	err = r.Get(ctx, client.ObjectKey{Name: cronJobName, Namespace: dump.Namespace}, cronJob)
	if err != nil {
		if apierrors.IsNotFound(err) {
			// CronJob doesn't exist, create it
			cronJob = r.constructCronJob(dump)
			if err := controllerutil.SetControllerReference(dump, cronJob, r.Scheme); err != nil {
				log.Error(err, "Failed to set controller reference")
				return ctrl.Result{}, err
			}
			if err := r.Create(ctx, cronJob); err != nil {
				log.Error(err, "Failed to create CronJob", "cronjob", cronJobName)
				return ctrl.Result{}, err
			}
			log.Info("Successfully created CronJob", "cronjob", cronJobName)
		} else {
			log.Error(err, "Failed to get CronJob", "cronjob", cronJobName)
			return ctrl.Result{}, err
		}
	} else {
		log.Info("CronJob already exists", "cronjob", cronJobName)
	}

	// Update status to track that we've observed this generation
	dump.Status.ObservedGeneration = dump.Generation
	if err := r.Status().Update(ctx, dump); err != nil {
		log.Error(err, "Failed to update Dump status")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// constructCronJob creates a CronJob object for the given Dump resource
func (r *DumpReconciler) constructCronJob(dump *dumpv1alpha1.Dump) *batchv1.CronJob {
	cronJob := &batchv1.CronJob{
		ObjectMeta: metav1.ObjectMeta{
			Name:      dump.Name + "-cronjob",
			Namespace: dump.Namespace,
			Labels: map[string]string{
				"app":       "dump-cronjob",
				"dump-name": dump.Name,
			},
		},
		Spec: batchv1.CronJobSpec{
			Schedule: dump.Spec.Schedule,
			JobTemplate: batchv1.JobTemplateSpec{
				Spec: batchv1.JobSpec{
					Template: corev1.PodTemplateSpec{
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name:    "dump-container",
									Image:   "busybox:latest",
									Command: []string{"/bin/sh", "-c", "echo 'Running dump job for " + dump.Spec.TcpFilter + "'"},
								},
							},
							RestartPolicy: corev1.RestartPolicyOnFailure,
						},
					},
				},
			},
		},
	}
	return cronJob
}

// SetupWithManager sets up the controller with the Manager.
func (r *DumpReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&dumpv1alpha1.Dump{}).
		Named("dump").
		Complete(r)
}
