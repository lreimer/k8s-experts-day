/*
Copyright 2021 Mario-Leander Reimer.
*/

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// DemoSpec defines the desired state of Demo
type DemoSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of Demo. Edit Demo_types.go to remove/update
	Foo string `json:"foo,omitempty"`

	// Message to the World
	Message string `json:"message,omitempty"`
}

// DemoStatus defines the observed state of Demo
type DemoStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Demo is the Schema for the demoes API
type Demo struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DemoSpec   `json:"spec,omitempty"`
	Status DemoStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// DemoList contains a list of Demo
type DemoList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Demo `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Demo{}, &DemoList{})
}
