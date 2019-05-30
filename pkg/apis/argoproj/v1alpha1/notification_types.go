/*

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
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// NotificationSpec defines the desired state of Notification
type NotificationSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	NotificationConfig `json:",inline"`
}

const (
	RULES_OPERATOR_EQ = "eq"
	RULES_OPERATOR_NE = "ne"
	RULES_OPERATOR_GT = "gt"
	RULES_OPERATOR_GE = "ge"
	RULES_OPERATOR_LT = "lt"
	RULES_OPERATOR_LE = "le"
)

const (
	RULES_LOGICAL_AND = "and"
	RULES_LOGICAL_OR  = "or"
)

const (
	NOTIFICATION_LEVEL_INFO     = "info"
	NOTIFICATION_LEVEL_WARNING  = "warning"
	NOTIFICATION_LEVEL_CRITICAL = "critical"
)

// NotificationStatus defines the observed state of Notification
type NotificationStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	RuleStatus []RuleStatus `json:",rules,omitempty"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Notification is the Schema for the notifications API
// +k8s:openapi-gen=true
type Notification struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              NotificationConfig `json:"spec,omitempty"`
	Status            NotificationStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// NotificationList contains a list of Notification
type NotificationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Notification `json:"items"`
}

type NotificationConfig struct {
	KubeConfigPath  string                      `json:"kubeConfigPath,omitempty"`
	Namespace       string                      `json:"namespace,omitempty"`
	MonitorResource schema.GroupVersionResource `json:"monitorResource,omitempty"`
	Rules           []Rule                      `json:"rules,omitempty"`
	Notifier        []Notifier                  `json:"notifiers,omitempty"`
}

type Notifier struct {
	Name  string        `json:"name"`
	Slack SlackNotifier `json:"slack,omitempty"`
	Email EmailNotifier `json:"email,omitempty"`
}

type SlackNotifier struct {
	HookUrlSecret apiv1.SecretKeySelector `json:"hookUrlSecret,omitempty"`
	TokenSecret   apiv1.SecretKeySelector `json:"tokenSecret,omitempty"`
	Channel       string                  `json:"channel,omitempty"`
}

type EmailNotifier struct {
	SmtpHost       string                  `json:"smtphost,omitempty"`
	SmtpPort       int                     `json:"smtpport,omitempty"`
	UserNameSecret apiv1.SecretKeySelector `json:"usernameSecret"`
	PasswordSecret apiv1.SecretKeySelector `json:"passwordSecret"`
	FromEmailId    string                  `json:"fromEmailId,omitempty"`
	SenderList     []string                `json:"senderList,omitempty"`
}

type Rule struct {
	Name            string      `json:"name"`
	InitialDelaySec int         `json:"initialDelaySec"`
	ThrottleMinutes int         `json:"throttleMintues"`
	AllConditions   []Condition `json:"allConditions,omitempty"`
	AnyConditions   []Condition `json:"anyConditions,omitempty"`
	Events          []Event     `json:"events,omitempty"`
}

type RuleStatus struct {
	Name           string          `json:"name"`
	TriggeredCount int             `json:"triggeredCount"`
	ActiveTriggers []ActiveTrigger `json:"activeTriggers"`
}

type ActiveTrigger struct {
	Name          string      `json:"name"`
	LastTriggered metav1.Time `json:"lastTriggered"`
}
type Condition struct {
	Jsonpath        string      `json:"jsonPath,omitempty"`
	Operator        string      `json:"operator,omitempty"`
	Value           string      `json:"value,omitempty"`
	ValueJsonPath   string      `json:"valueJsonPath,omitempty"`
	JoinOperator    string      `json:"joinOperator,omitempty"`
	ChildConditions []Condition `json:"childConditions,omitempty"`
}
type Event struct {
	Message           string   `json:"message,omitempty"`
	EmailSubject      string   `json:"emailSubject,omitempty"`
	NotificationLevel string   `json:"notificationLevel,omitempty"`
	NotifierNames     []string `json:"notifierNames,omitempty"`
}

func init() {
	SchemeBuilder.Register(&Notification{}, &NotificationList{})
}
