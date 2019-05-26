package config

//type NotificationConfig struct {
//	Notifications []Notification `json:"notifications,omitempty"`
//}
//
//type Notification struct {
//	Name            string                       `json:"name"`
//	TriggerAction   string                       `json:"triggerAction"`
//	KubeConfigPath  string                       `json:"kubeConfigPath,omitempty"`
//	Namespace       string                       `json:"namespace,omitempty"`
//	MonitorResource *schema.GroupVersionResource `json:"monitorResource,omitempty"`
//	Rules           []Rule                       `json:"rules,omitempty"`
//	Notifier        []Notifier                   `json:"notifiers,omitempty"`
//}
//
//type Notifier struct {
//	Name  string         `json:"name"`
//	Slack *SlackNotifier `json:"slack,omitempty"`
//	Email *EmailNotifier `json:"email,omitempty"`
//}
//
//type SlackNotifier struct {
//	HookUrl string `json:"hookurl,omitempty"`
//	Token   string `json:"token,omitempty"`
//	Channel string `json:"channel,omitempty"`
//}
//
//type EmailNotifier struct {
//	SmtpHost       string                  `json:"smtphost,omitempty"`
//	SmtpPort       int                     `json:"smtpport,omitempty"`
//	UserNameSecret apiv1.SecretKeySelector `json:"usernameSecret"`
//	PasswordSecret apiv1.SecretKeySelector `json:"passwordSecret"`
//	FromEmailId    string                  `json:"fromEmailId,omitempty"`
//	SenderList     []string                `json:"senderList,omitempty"`
//}
//
//type Rule struct {
//	AllConditions []Condition `json:"allConditions,omitempty"`
//	AnyConditions []Condition `json:"anyConditions,omitempty"`
//	Events        []Event     `json:"events,omitempty"`
//}
//type Condition struct {
//	Jsonpath        string      `json:"jsonPath,omitempty"`
//	Operator        string      `json:"operator,omitempty"`
//	Value           string      `json:"value,omitempty"`
//	ValueJsonPath   string      `json:"valueJsonPath,omitempty"`
//	JoinOperator    string      `json:"joinOperator,omitempty"`
//	ChildConditions []Condition `json:"childconditions,omitempty"`
//}
//type Event struct {
//	Message           string   `json:"message,omitempty"`
//	EmailSubject      string   `json:"emailSubject,omitempty"`
//	NotificationLevel string   `json:"notificationLevel,omitempty"`
//	NotifierNames     []string `json:"notifierNames,omitempty"`
//}
