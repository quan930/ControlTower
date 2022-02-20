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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// HookSpec defines the desired state of Hook
type HookSpec struct {
	Hooks       []HookItem   `json:"hooks"`
	GitEvents   []GitEvent   `json:"git_events"`
	ImageEvents []ImageEvent `json:"image_events"`
}

// HookStatus defines the observed state of Hook
type HookStatus struct {
	GitEventHistory   []GitEventHistory `json:"git_event_history"`
	ImageEventHistory []ImageEvent      `json:"image_event_history"`
}

//+kubebuilder:object:root=true

// Hook is the Schema for the hooks API
//+kubebuilder:subresource:status
type Hook struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   HookSpec   `json:"spec,omitempty"`
	Status HookStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// HookList contains a list of Hook
type HookList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Hook `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Hook{}, &HookList{})
}

//HookItem todo 校验
type HookItem struct {
	GitRepository     string   `json:"git_repository"`
	Branches          []string `json:"branches"`
	ImageRepository   string   `json:"image_repository"`
	Dockerfile        string   `json:"dockerfile"`
	ImageRepoUser     string   `json:"image_repo_user"`
	ImageRepoPassword string   `json:"image_repo_password"`
	ImageBuild        bool     `json:"image_build"`
	UpdateImage       bool     `json:"update_image"`
}

//GitEvent todo 校验
type GitEvent struct {
	GitRepository string `json:"git_repository"`
	Branch        string `json:"branch"`
}

//GitEventHistory git event history
type GitEventHistory struct {
	GitRepository string `json:"git_repository"`
	Branch        string `json:"branch"`
	DateTime      string `json:"date_time"`
	Status        string `json:"status"`
}

//ImageEvent todo 校验
type ImageEvent struct {
	ImageRepository string `json:"image_repository"`
}
