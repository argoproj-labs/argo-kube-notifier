package integration

import (
	"fmt"
	"github.com/argoproj-labs/argo-kube-notifier/pkg/apis/argoproj/v1alpha1"
	"github.com/nlopes/slack"
)

type SlackClient struct {
	slackConfig *v1alpha1.SlackNotifier
}

func NewSlackClient(slackConfig *v1alpha1.SlackNotifier) *SlackClient {
	slackClient := SlackClient{slackConfig: slackConfig}
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
	err := slack.PostWebhook(s.slackConfig.HookUrl, &slack.WebhookMessage{
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
	err := slack.PostWebhook(s.slackConfig.HookUrl, &slack.WebhookMessage{
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
	err := slack.PostWebhook(s.slackConfig.HookUrl, &slack.WebhookMessage{Attachments: []slack.Attachment{attachment}, Channel: s.slackConfig.Channel, IconEmoji: "+1"})
	if err != nil {
		fmt.Println(err)
	}
	return nil
}
func (s *SlackClient) SendInfoNotification(msg ...string) error {
	return nil
}
