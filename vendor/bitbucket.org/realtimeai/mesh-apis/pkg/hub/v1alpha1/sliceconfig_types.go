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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// SliceConfigSpec defines the desired state of SliceConfig
type SliceConfigSpec struct {
	SliceSubnet string `json:"sliceSubnet,omitempty"`
	//+kubebuilder:validation:Enum:=Application
	SliceType                 string                    `json:"sliceType,omitempty"`
	SpokeSliceGatewayProvider SpokeSliceGatewayProvider `json:"spokeSliceGatewayProvider,omitempty"`
	//+kubebuilder:validation:Enum:=Local
	SliceIpamType             string                    `json:"sliceIpamType,omitempty"`
	Clusters                  []string                  `json:"clusters,omitempty"`
	SchSliceClusterSite       string                    `json:"schSliceClusterSite,omitempty"` //FIXME:: what's the purpose
	QosProfileDetails         QOSProfile                `json:"qosProfileDetails,omitempty"`
	SliceAdmissionIndex       int                       `json:"sliceAdmissionIndex,omitempty"` //FIXME:: need discussion in nodeport
	NamespaceIsolationProfile NamespaceIsolationProfile `json:"namespaceIsolationProfile,omitempty"`
	NetworkPolicyVersion      int                       `json:"networkPolicyVersion,omitempty"` //FIXME::
}

// SpokeSliceGatewayProvider defines the configuration for slicegateway
type SpokeSliceGatewayProvider struct {
	//+kubebuilder:validation:Enum:=OpenVPN
	SpokeSliceGatewayType string `json:"spokeSliceGatewayType,omitempty"`
	//+kubebuilder:validation:Enum:=Local
	SliceCAType string `json:"sliceCaType,omitempty"`
	//+kubebuilder:validation:Enum:=IDPCognito;IDPOkta;IDPOAuth;IDPOpenID
	SliceIdp          string `json:"sliceIdp,omitempty"`
	MaxGateways       int    `json:"maxGateways,omitempty"`
	GatewaySubnetSize int    `json:"gatewaySubnetSize,omitempty"`
	GatewayQosProfile string `json:"gatewayQosProfile,omitempty"`
}

// QOSProfile is the QOS Profile configuration from backend
type QOSProfile struct {
	//+kubebuilder:validation:Enum:=HTB
	QueueType               string `json:"queueType,omitempty"`
	Priority                string `json:"priority,omitempty"`
	TcType                  string `json:"tcType,omitempty"`
	BandwidthCeilingKbps    int    `json:"bandwidthCeilingKbps,omitempty"` //FIXME:: For unlimited
	BandwidthGuaranteedKbps int    `json:"bandwidthGuaranteedKbps,omitempty"`
	//+kubebuilder:validation:Enum:=Default;AF11;AF12;AF13;AF21;AF22;AF23;AF31;AF32;AF33;AF41;AF42;AF43;EF
	DscpClass string `json:"dscpClass,omitempty"`
}

type NamespaceIsolationProfile struct {
	IsolationEnabled          bool                      `json:"isolationEnabled,omitempty"`
	ApplicationNamespaces     []string                  `json:"applicationNamespaces,omitempty"` //FIXME::
	AllowedNamespaces         []string                  `json:"allowedNamespaces,omitempty"`     //FIXME::
	SliceNamespaceHierarchy   []SliceNamespaceHierarchy `json:"sliceNamespaceHierarchy,omitempty"`
	NamespaceHierarchyEnabled bool                      `json:"namespaceHierarchyEnabled,omitempty"`
}

type SliceNamespaceHierarchy struct {
	Cluster       []string        `json:"cluster,omitempty"`
	Namespace     string          `json:"namespace,omitempty"`
	SubNamespaces []SubNamespaces `json:"subNamespaces,omitempty"`
}

//+kubebuilder:pruning:PreserveUnknownFields
//+kubebuilder:validation:EmbeddedResource
type SubNamespaces unstructured.Unstructured

func (in *SubNamespaces) DeepCopyInto(out *SubNamespaces) {
	// controller-gen cannot handle the interface{} type of an aliased Unstructured, thus we write our own DeepCopyInto function.
	if out != nil {
		casted := unstructured.Unstructured(*in)
		deepCopy := casted.DeepCopy()
		out.Object = deepCopy.Object
	}
}

// SliceConfigStatus defines the observed state of SliceConfig
type SliceConfigStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// SliceConfig is the Schema for the sliceconfig API
type SliceConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SliceConfigSpec   `json:"spec,omitempty"`
	Status SliceConfigStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// SliceConfigList contains a list of SliceConfig
type SliceConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SliceConfig `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SliceConfig{}, &SliceConfigList{})
}
