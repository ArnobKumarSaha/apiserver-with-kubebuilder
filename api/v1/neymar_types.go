/*
Copyright 2021.

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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// NeymarSpec defines the desired state of Neymar
type NeymarSpec struct {
	DeploymentName string `json:"deploymentName"`
	DeploymentImage string `json:"deploymentImage"`
	Replicas       *int32 `json:"replicas"`
	ServiceName 	string `json:"serviceName"`
	ServicePort 	int32 `json:"servicePort"`
	ServiceType     string `json:"serviceType"`
	ServiceTargetPort int32 `json:"serviceTargetPort"`
}

// NeymarStatus defines the observed state of Neymar
type NeymarStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	AvailableReplicas int32 `json:"availableReplicas"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Neymar is the Schema for the neymars API
type Neymar struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   NeymarSpec   `json:"spec,omitempty"`
	Status NeymarStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// NeymarList contains a list of Neymar
type NeymarList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Neymar `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Neymar{}, &NeymarList{})
}
