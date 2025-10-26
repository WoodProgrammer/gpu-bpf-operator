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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ProbeTargetBindingSpec defines the desired state of ProbeTargetBinding.
type ProbeTargetBindingSpec struct {
	PolicyRef      string            `json:"policyRef"`
	NodeSelector   map[string]string `json:"nodeSelector,omitempty"`
	CanaryPercent  int               `json:"canaryPercent,omitempty"`
	MaxUnavailable int               `json:"maxUnavailable,omitempty"`
}
type ProbeTargetBindingStatus struct {
	AppliedHash string `json:"appliedHash,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// ProbeTargetBinding is the Schema for the probetargetbindings API.
type ProbeTargetBinding struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ProbeTargetBindingSpec   `json:"spec,omitempty"`
	Status ProbeTargetBindingStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ProbeTargetBindingList contains a list of ProbeTargetBinding.
type ProbeTargetBindingList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ProbeTargetBinding `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ProbeTargetBinding{}, &ProbeTargetBindingList{})
}
