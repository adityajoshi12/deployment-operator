/*
Copyright 2023.

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

// ScaleSpec defines the desired state of Scale
type ScaleSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	StartTime      int    `json:"startTime"`
	EndTime        int    `json:"endTime"`
	MinReplicas    int    `json:"minReplicas"`
	MaxReplicas    int    `json:"maxReplicas"`
	DeploymentName string `json:"deploymentName"`
	// +kubebuilder:default=default
	DeploymentNamespace string `json:"deploymentNamespace"`
}

// ScaleStatus defines the observed state of Scale
type ScaleStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Scale is the Schema for the scales API
type Scale struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ScaleSpec   `json:"spec,omitempty"`
	Status ScaleStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ScaleList contains a list of Scale
type ScaleList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Scale `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Scale{}, &ScaleList{})
}
