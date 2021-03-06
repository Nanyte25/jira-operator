// Copyright 2018 Jira Operator Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package v1alpha1

import (
	"fmt"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// DefaultBaseImage is the default docker image to use for JIRA Pods.
	DefaultBaseImage = "cptactionhank/atlassian-jira-software"

	// DefaultVersion is the default Jira version to use for JIRA Pods.
	DefaultVersion = "7.10.2"

	// DefaultDataMountPath is the default filesystem path for JIRA Home.
	DefaultDataMountPath = "/var/atlassian/jira"

	// DefaultIngressPath is the default path for the Ingress resource.
	DefaultIngressPath = "/"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// JiraList resource
type JiraList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []Jira `json:"items"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Jira resource
type Jira struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Spec              JiraSpec   `json:"spec"`
	Status            JiraStatus `json:"status,omitempty"`
}

// JiraPodPolicy defines the policy for pods owned by JIRA operator.
type JiraPodPolicy struct {
	// Resources is the resource requirements for the jira container.
	Resources v1.ResourceRequirements `json:"resources,omitempty"`

	// PersistentVolumeClaimSpec is the spec to describe PVC for the jira container
	// This field is optional. If no PVC spec, jira container will use emptyDir as volume
	PersistentVolumeClaimSpec *v1.PersistentVolumeClaimSpec `json:"persistentVolumeClaimSpec,omitempty"`
}

// JiraIngressPolicy defines the Ingress policy for the operator.
type JiraIngressPolicy struct {
	Annotations map[string]string `json:"annotations,omitempty"`
	Host        string            `json:"host"`
	TLS         bool              `json:"tls"`
	Path        string            `json:"path"`
	SecretName  string            `json:"secretName,omitempty"`
}

// JiraSpec resource
type JiraSpec struct {
	// BaseImage image to use for a JIRA deployment.
	BaseImage string `json:"baseImage"`

	// Version is the Jira version to deploy.
	Version string `json:"version"`

	// DataMountPath path for JIRA Home.
	DataMountPath string `json:"dataMountPath"`

	// ConfigMapName is the name of ConfigMap to use or create.
	ConfigMapName string `json:"configMapName"`

	// SecretName is the name of Secret to use or create.
	SecretName string `json:"secretName"`

	// Pod defines the policy for pods owned by JIRA operator.
	// This field cannot be updated once the CR is created.
	Pod *JiraPodPolicy `json:"pod,omitempty"`

	// Ingress defines the Ingress policy for the JIRA operator.
	Ingress *JiraIngressPolicy `json:"ingress,omitempty"`
}

// SetDefaults sets the default vaules for the Jira spec and returns true if the spec was changed
func (j *Jira) SetDefaults() bool {
	changed := false
	if len(j.Spec.BaseImage) == 0 {
		j.Spec.BaseImage = DefaultBaseImage
		changed = true
	}
	if len(j.Spec.Version) == 0 {
		j.Spec.Version = DefaultVersion
		changed = true
	}
	if len(j.Spec.ConfigMapName) == 0 {
		j.Spec.ConfigMapName = j.ObjectMeta.Name
		changed = true
	}
	if len(j.Spec.DataMountPath) == 0 {
		j.Spec.DataMountPath = DefaultDataMountPath
		changed = true
	}
	if len(j.Spec.SecretName) == 0 {
		j.Spec.SecretName = j.ObjectMeta.Name
		changed = true
	}
	if j.SetIngressDefaults() {
		changed = true
	}
	return changed
}

// SetIngressDefaults sets the default vaules for the Jira Ingress policy and returns true if the spec was changed.
func (j *Jira) SetIngressDefaults() bool {
	changed := false
	if j.Spec.Ingress == nil {
		return changed
	}
	if len(j.Spec.Ingress.Host) == 0 {
		j.Spec.Ingress.Host = j.ObjectMeta.Name
	}
	if len(j.Spec.Ingress.Path) == 0 {
		j.Spec.Ingress.Path = DefaultIngressPath
	}
	if len(j.Spec.Ingress.SecretName) == 0 {
		j.Spec.Ingress.SecretName = fmt.Sprintf("%s-ingress", j.ObjectMeta.Name)
		changed = true
	}
	return changed
}

// IsIngressEnabled shortcut fucntion to determine Ingress status.
func (j *Jira) IsIngressEnabled() bool {
	return j.Spec.Ingress != nil
}

// IsIngressTLSEnabled shortcut fucntion to determine Ingress TLS status.
func (j *Jira) IsIngressTLSEnabled() bool {
	if j.IsIngressEnabled() {
		return j.Spec.Ingress.TLS
	}
	return false
}

// IsPodPolicySet shortcut fucntion to determine Pod policy status.
func (j *Jira) IsPodPolicySet() bool {
	return j.Spec.Pod != nil
}

// IsPVEnabled shortcut fucntion to determine PV status.
func (j *Jira) IsPVEnabled() bool {
	if j.IsPodPolicySet() {
		return j.Spec.Pod.PersistentVolumeClaimSpec != nil
	}
	return false
}

// JiraStatus resource
type JiraStatus struct {
	// Endpoint is the URI for accessing Jira.
	Endpoint string `json:"endpoint,omitempty"`

	// ServiceName is the LB service for accessing Jira.
	ServiceName string `json:"serviceName,omitempty"`

	// State is the current state for the Jira application (Initializing, Running, etc...).
	State string `json:"state"`
}
