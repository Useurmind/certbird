package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CertRequest is 
type CertRequest struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`

	Spec   CertRequestSpec   `json:"spec"`
	Status CertRequestStatus `json:"status,omitempty"`
}

// CertRequestSpec is the spec for a CertRequest resource
type CertRequestSpec struct {
	Foo string `json:"foo"`
	Bar bool   `json:"bar"`
}

// CertRequestStatus is the status for a CertRequest resource
type CertRequestStatus struct {
	State   CertRequestState `json:"state,omitempty"`
	Message string       `json:"message,omitempty"`
}

type CertRequestState string

const (
	ExampleStateCreated   CertRequestState = "Created"
	ExampleStateProcessed CertRequestState = "Processed"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CertRequestList is a list of CertRequest resources
type CertRequestList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []CertRequest `json:"items"`
}
