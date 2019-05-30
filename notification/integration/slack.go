package integration

import (
	"fmt"
	"github.com/argoproj-labs/argo-kube-notifier/pkg/apis/argoproj/v1alpha1"
	"github.com/argoproj-labs/argo-kube-notifier/util"
	"github.com/nlopes/slack"
	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
)

type SlackClient struct {
	slackConfig *v1alpha1.SlackNotifier
	hookUrl     string
}

func NewSlackClient(clientSet kubernetes.Interface, namespace string, slackConfig *v1alpha1.SlackNotifier) *SlackClient {
	slackClient := SlackClient{slackConfig: slackConfig}

	hookUrlBytes, err := util.GetSecrets(clientSet, namespace, slackConfig.HookUrlSecret.Name, slackConfig.HookUrlSecret.Key)
	if err != nil {
		log.Warnf("Slack hook url secret failed to read. %v", err)
	}
	slackClient.hookUrl = string(hookUrlBytes)

	return &slackClient
}

//err := slack.PostWebhook("https://hooks.slack.com/services/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", &slack.WebhookMessage{ Text : msg } )

func (s *SlackClient) SendSuccessNotification(msg ...string) error {
	attachment := slack.Attachment{
		Color:    "good",
		Text:     msg[0],
		Title:    "Argo Kube Notifier",
		ThumbURL: "http://chittagongit.com/images/icon-success/icon-success-17.jpg",
	}
	err := slack.PostWebhook(s.hookUrl, &slack.WebhookMessage{
		Attachments: []slack.Attachment{attachment}, Channel: s.slackConfig.Channel, IconEmoji: "+1"})
	if err != nil {
		fmt.Println(err)
	}
	return nil
}

func (s *SlackClient) SendWarningNotification(msg ...string) error {
	attachment := slack.Attachment{
		Color:    "warning",
		Text:     msg[0],
		Title:    "Argo Kube Notifier",
		ThumbURL: "https://icon2.kisspng.com/20180626/kiy/kisspng-warning-sign-computer-icons-clip-art-warning-icon-5b31bd67368be5.4827407215299864072234.jpg",
	}
	err := slack.PostWebhook(s.hookUrl, &slack.WebhookMessage{
		Attachments: []slack.Attachment{attachment}, Channel: s.slackConfig.Channel, IconEmoji: "+1"})
	if err != nil {
		fmt.Println(err)
	}
	return nil
}
func (s *SlackClient) SendFailledNotification(msg ...string) error {
	attachment := slack.Attachment{
		Color:    "danger",
		Text:     msg[0],
		Title:    "Argo Kube Notifier",
		ThumbURL: "https://cdn3.iconfinder.com/data/icons/picons-weather/57/53_warning-512.png",
	}
	err := slack.PostWebhook(s.hookUrl, &slack.WebhookMessage{Attachments: []slack.Attachment{attachment}, Channel: s.slackConfig.Channel, IconEmoji: "+1"})
	if err != nil {
		fmt.Println(err)
	}
	return nil
}
func (s *SlackClient) SendInfoNotification(msg ...string) error {
	return nil
}
