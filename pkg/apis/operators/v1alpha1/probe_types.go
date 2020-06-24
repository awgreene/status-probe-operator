package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ProbeSpec defines the desired state of Probe
type ProbeSpec struct {
	// specOverride allows annotations to override the spec. This feature may be expanded to accept other overrides in the future.
	// +kubebuilder:validation:Enum=crdAnnotations
	SpecOverride   string          `json:"specOverride,omitempty"`
	ProbeResources []ProbeResource `json:"probeResources,omitempty"`
}

// ProbeResource
type ProbeResource struct {
	// name is the name of the CRD used for this probe.
	Name string `json:"crdName"`

	// Upgradeable is a list of condition types that map to the upgradeable "OLM Supported Condition"
	Upgradeable []string `json:"upgradeable,omitempty"`

	// Important is a list of condition types that map to the important "OLM Supported Condition"
	Important []string `json:"important,omitempty"`
}

type ProbeResourceStatus struct {
	Name      string                   `json:"name"`
	Resources []ProbeResourceCondition `json:"resources,omitempty"`
}

type ProbeResourceCondition struct {
	UID types.UID `json:"uid"`
	// conditions
	Conditions []ProbeCondition `json:"conditions,omitempty"`
}

// ProbeStatus defines the observed state of Probe
type ProbeStatus struct {
	// upgradeable communicates to OLM if the operator can be upgraded. If the field is not set, OLM will assume that the operator has not communicated this state.
	Upgradeable *bool `json:"upgradeable,omitempty"`

	// resource is a list of resources.
	ProbeResources []ProbeResourceStatus `json:"probeResources,omitempty"`
}

type ProbeCondition struct {
	Type               string `json:"type,omitempty"`
	Reason             string `json:"reason,omitempty"`
	Message            string `json:"message,omitempty"`
	LastTransitionTime string `json:"lastTransitionTime,omitempty"`
	Status             string `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Probe is the Schema for the probes API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=probes,scope=Namespaced
type Probe struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ProbeSpec   `json:"spec,omitempty"`
	Status ProbeStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ProbeList contains a list of Probe
type ProbeList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Probe `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Probe{}, &ProbeList{})
}
