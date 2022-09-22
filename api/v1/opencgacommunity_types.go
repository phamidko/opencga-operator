/*
Copyright 2022.

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
	"encoding/json"

	"github.com/stretchr/objx"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type Type string

const (
	ReplicaSet Type = "ReplicaSet"
)

type Phase string

const (
	Running Phase = "Running"
	Failed  Phase = "Failed"
	Pending Phase = "Pending"
)

const (
	defaultClusterDomain = "cluster.local"
)

// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// OpenCGACommunitySpec defines the desired state of OpenCGACommunity
type OpenCGACommunitySpec struct {
	// Members is the number of members in the replica set
	Members int `json:"members"`
	// Type defines which type of OpenCGA REST deployment the resource should create
	// +kubebuilder:validation:Enum=ReplicaSet
	Type Type `json:"type"`
	// Version defines which version of OpenCGA will be used
	Version string `json:"version"`
}

// OpenCGAConfiguration holds the optional openCGA REST configuration
// that should be merged with the operator created one.
//
// The CRD generator does not support map[string]interface{}
// on the top level and hence we need to work around this with
// a wrapping struct.
type OpenCGAConfiguration struct {
	Object map[string]interface{} `json:"-"`
}

// MarshalJSON defers JSON encoding to the wrapped map
func (m *OpenCGAConfiguration) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.Object)
}

// UnmarshalJSON will decode the data into the wrapped map
func (m *OpenCGAConfiguration) UnmarshalJSON(data []byte) error {
	if m.Object == nil {
		m.Object = map[string]interface{}{}
	}

	return json.Unmarshal(data, &m.Object)
}

func (m *OpenCGAConfiguration) DeepCopy() *OpenCGAConfiguration {
	return &OpenCGAConfiguration{
		Object: runtime.DeepCopyJSON(m.Object),
	}
}

// NewOpenCGAConfiguration returns an empty NewOpenCGAConfiguration
func NewOpenCGAConfiguration() OpenCGAConfiguration {
	return OpenCGAConfiguration{Object: map[string]interface{}{}}
}

// SetOption updated the OpenCGAConfiguration with a new option
func (m OpenCGAConfiguration) SetOption(key string, value interface{}) OpenCGAConfiguration {
	m.Object = objx.New(m.Object).Set(key, value)
	return m
}

// OpenCGACommunityStatus defines the observed state of OpenCGACommunity
type OpenCGACommunityStatus struct {
	RestURI string `json:"opencgarestUri"`
	Phase   Phase  `json:"phase"`
	Version string `json:"version"`

	CurrentStatefulSetReplicas int `json:"currentStatefulSetReplicas"`
	CurrentRestMembers         int `json:"currentOpenCGARESTMembers"`

	Message string `json:"message,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// OpenCGACommunity is the Schema for the opencgacommunities API
// +kubebuilder:resource:path=opencgacommunity,scope=Namespaced,shortName=ocbc,singular=opencgacommunity
// +kubebuilder:printcolumn:name="RestURI",type="string",JSONPath=".status.opencgarestUri",description="Current REST URI of the OpenCGA REST deployment"
// +kubebuilder:printcolumn:name="Phase",type="string",JSONPath=".status.phase",description="Current state of the OpenCGA REST deployment"
// +kubebuilder:printcolumn:name="Version",type="string",JSONPath=".status.version",description="Version of OpenCGA REST server"
type OpenCGACommunity struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   OpenCGACommunitySpec   `json:"spec,omitempty"`
	Status OpenCGACommunityStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// OpenCGACommunityList contains a list of OpenCGACommunity
type OpenCGACommunityList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []OpenCGACommunity `json:"items"`
}

func init() {
	SchemeBuilder.Register(&OpenCGACommunity{}, &OpenCGACommunityList{})
}
