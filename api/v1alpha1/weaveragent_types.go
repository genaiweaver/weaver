// SPDX-License-Identifier: Apache-2.0
// +kubebuilder:object:generate=true
// +groupName=weaver.io

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// WeaverAgentSpec defines the desired state of WeaverAgent.
// It declares a high-level, declarative graph of Knative Services and Triggers.
// +kubebuilder:validation:Required
type WeaverAgentSpec struct {
	// Broker is the name of the Knative Broker through which all events flow.
	// Defaults to "default" if not specified.
	// +kubebuilder:validation:MinLength=1
	Broker string `json:"broker"`

	// Nodes is the list of service nodes participating in this agent graph.
	// Each node corresponds to a Knative Service name.
	// +kubebuilder:validation:MinItems=1
	Nodes []AgentNode `json:"nodes"`

	// Edges defines the directed event flow between nodes.
	// Each edge creates a Knative Trigger on the specified Broker.
	// +kubebuilder:validation:MinItems=1
	Edges []AgentEdge `json:"edges"`
}

// AgentNode represents a single node in the agent graph.
// It maps directly to a Knative Service.
type AgentNode struct {
	// ServiceName is the Knative Service name for this node.
	// +kubebuilder:validation:MinLength=1
	ServiceName string `json:"serviceName"`
}

// AgentEdge defines an event-driven connection between two nodes.
// The controller creates a Knative Trigger for each edge.
type AgentEdge struct {
	// EventType is the CloudEvent 'type' attribute to filter on.
	// When the Broker receives an event of this type, it will route to the 'To' service.
	// +kubebuilder:validation:MinLength=1
	EventType string `json:"eventType"`

	// To is the ServiceName of the destination node for this edge.
	// Must match one of the names listed under spec.nodes.
	// +kubebuilder:validation:MinLength=1
	To string `json:"to"`

	// Filter allows specifying additional CloudEvent attribute filters
	// (key = attribute name, value = exact match).
	// +optional
	Filter map[string]string `json:"filter,omitempty"`
}

// WeaverAgentStatus defines the observed state of WeaverAgent.
// It reflects the controller’s last reconciliation.
type WeaverAgentStatus struct {
	// ObservedGeneration is the last reconciled generation.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`

	// Conditions represent the latest observations of the agent’s state.
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:resource:shortName=wa;wagent
//+kubebuilder:subresource:status

// WeaverAgent is the Schema for the weaveragents API.
// It encapsulates the full agent workflow graph in one resource.
type WeaverAgent struct {
	metav1.TypeMeta   `json:",inline"`            // API version & kind
	metav1.ObjectMeta `json:"metadata,omitempty"` // Name, namespace, labels, annotations

	Spec   WeaverAgentSpec   `json:"spec,omitempty"`   // Desired state (graph definition)
	Status WeaverAgentStatus `json:"status,omitempty"` // Observed state (health & generation)
}

//+kubebuilder:object:root=true

// WeaverAgentList contains a list of WeaverAgent resources.
type WeaverAgentList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []WeaverAgent `json:"items"`
}

func init() {
	SchemeBuilder.Register(&WeaverAgent{}, &WeaverAgentList{})
}
