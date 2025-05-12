// api/v1alpha1/weavernode_types.go
// +k8s:deepcopy-gen=package
package v1alpha1

import (
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// WeaverNodeSpec defines the desired state of WeaverNode
type WeaverNodeSpec struct {
    ID                   string            `json:"id"`
    ComponentType        string            `json:"componentType"`
    ServiceName          string            `json:"serviceName"`
    Endpoint             string            `json:"endpoint"`
    Consumes             []string          `json:"consumes,omitempty"`
    Produces             []string          `json:"produces,omitempty"`
    InvocationMechanisms []InvocationSpec  `json:"invocationMechanisms,omitempty"`
    EventMappings        []EventMapping    `json:"eventMappings,omitempty"`
    Metadata             map[string]string `json:"metadata,omitempty"`
    StoreInRedis         bool              `json:"storeInRedis,omitempty"`
    RedisConfigRef       *RedisConfig      `json:"redisConfigRef,omitempty"`
}

// InvocationSpec declares how to invoke the node
type InvocationSpec struct {
    Type       string            `json:"type"`                 // e.g. "cloudEvent" or "http"
    Broker     string            `json:"broker,omitempty"`     // for cloudEvent
    Attributes map[string]string `json:"attributes,omitempty"` // ce-type, ce-source, etc.
    Method     string            `json:"method,omitempty"`     // for HTTP
    Path       string            `json:"path,omitempty"`
}

// EventMapping follows an OpenAI-style function schema for each event
type EventMapping struct {
    Name        string            `json:"name"`
    Description string            `json:"description,omitempty"`
    Parameters  map[string]string `json:"parameters,omitempty"` // JSON Schema as string map
}

// RedisConfig holds connection info, if caching to Redis
type RedisConfig struct {
    Host   string `json:"host"`
    Port   int    `json:"port"`
    Secret string `json:"secret,omitempty"`
}

// WeaverNodeStatus defines the observed state of WeaverNode
type WeaverNodeStatus struct {
    Healthy       bool              `json:"healthy,omitempty"`
    LastHeartbeat metav1.Time       `json:"lastHeartbeat,omitempty"`
    Conditions    []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// WeaverNode is the Schema for the weavernodes API
type WeaverNode struct {
    metav1.TypeMeta   `json:",inline"`
    metav1.ObjectMeta `json:"metadata,omitempty"`

    Spec   WeaverNodeSpec   `json:"spec,omitempty"`
    Status WeaverNodeStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// WeaverNodeList contains a list of WeaverNode
type WeaverNodeList struct {
    metav1.TypeMeta `json:",inline"`
    metav1.ListMeta `json:"metadata,omitempty"`
    Items           []WeaverNode `json:"items"`
}

func init() {
    SchemeBuilder.Register(&WeaverNode{}, &WeaverNodeList{})
}
