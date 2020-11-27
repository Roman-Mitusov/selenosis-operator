/*
Copyright 2020 Selenosis authors.

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
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// K8sSeleniumHubSpec defines the desired state of K8sSeleniumHub
type K8sSeleniumHubSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Selenosis replicas
	// +kubebuilder:validation:Minimum=1
	SelenosisReplicas int32 `json:"selenosisReplicas,omitempty"`
	// Selenoid ui port
	SelenoidUiPort int32 `json:"selenoidUiPort,omitempty"`
	// Selenosis port
	SelenosisPort int32 `json:"selenosisPort,omitempty"`
	// Browsers Proxy port
	BrowsersProxyPort int32 `json:"browsersProxyPort,omitempty"`
	// Selenosis image name
	SelenosisImage string `json:"selenosisImage"`
	// Selenoid UI image name
	SelenoidUiImage string `json:"selenoidUiImage"`
	// Selenoid UI Adapter image name
	SelenoidUiAdapterImage string `json:"selenoidUiAdapterImage"`
	// Browsers Proxy image name
	BrowsersProxyImage string `json:"browsersProxyImage"`
	// Image pull secret name. Use it only when you have images stored in private registry
	ImagePullSecretName string `json:"imagePullSecretName,omitempty"`
	// Service name fro browsers headless service
	BrowsersServiceName string `json:"browsersServiceName,omitempty"`
	// Service name for selenosis service
	SelenosisServiceName string `json:"selenosisServiceName,omitempty"`
	// Service name for selenoid ui service
	SelenoidUiServiceName string `json:"selenoidUiServiceName,omitempty"`
	// Browsers config definition
	BrowsersConfig BrowsersLayout `json:"browsersConfig"`
}

//Meta describes standart metadata
type Meta struct {
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
}

//Spec describes specification for Service
type Spec struct {
	Resources    apiv1.ResourceRequirements `json:"resources,omitempty"`
	HostAliases  []apiv1.HostAlias          `json:"hostAliases,omitempty"`
	EnvVars      []apiv1.EnvVar             `json:"envVars,omitempty"`
	NodeSelector map[string]string          `json:"nodeSelector,omitempty"`
	Affinity     apiv1.Affinity             `json:"affinity,omitempty"`
	DNSConfig    apiv1.PodDNSConfig         `json:"dnsConfig,omitempty"`
}

//Browsers Layout ...
type BrowsersLayout struct {
	DefaultSpec    Spec                   `json:"spec,omitempty"`
	Meta           Meta                   `json:"meta,omitempty"`
	Path           string                 `json:"path"`
	DefaultVersion string                 `json:"defaultVersion"`
	Versions       map[string]BrowserSpec `json:"versions"`
}

//BrowserSpec describes settings for Service
type BrowserSpec struct {
	BrowserName    string `json:"-"`
	BrowserVersion string `json:"-"`
	Image          string `json:"image"`
	Path           string `json:"path"`
	Meta           Meta   `json:"meta,omitempty"`
	Spec           Spec   `json:"spec,omitempty"`
}

// K8sSeleniumHubStatus defines the observed state of K8sSeleniumHub
type K8sSeleniumHubStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	SelenosisServiceReady        bool  `json:"selenosisServiceReady"`
	SelenoidUiServiceReady       bool  `json:"selenoidUiServiceReady"`
	HeadlessBrowsersServiceReady bool  `json:"headlessBrowsersServiceReady"`
	SelenosisReplicas            int32 `json:"selenosisReplicas"`
}

// +kubebuilder:object:root=true

// K8sSeleniumHub is the Schema for the k8sseleniumhubs API
type K8sSeleniumHub struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   K8sSeleniumHubSpec   `json:"spec,omitempty"`
	Status K8sSeleniumHubStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// K8sSeleniumHubList contains a list of K8sSeleniumHub
type K8sSeleniumHubList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []K8sSeleniumHub `json:"items"`
}

func init() {
	SchemeBuilder.Register(&K8sSeleniumHub{}, &K8sSeleniumHubList{})
}
