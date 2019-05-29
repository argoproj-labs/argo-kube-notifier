package integration

import (
	"github.com/argoproj-labs/argo-kube-notifier/pkg/apis/argoproj/v1alpha1"
	"gopkg.in/gomail.v2"
	"k8s.io/client-go/kubernetes"
)

type EmailClient struct {
	Client *gomail.Dialer
	Config *v1alpha1.EmailNotifier
}

func NewEmailClient(clientSet kubernetes.Interface, namespace string, emailConfig *v1alpha1.EmailNotifier) *EmailClient {
	emailClient := EmailClient{}
	//emailClient.Config = emailConfig
	//userNameBytes, err := util.GetSecrets(clientSet, namespace, emailConfig.UserNameSecret.Name, emailConfig.UserNameSecret.Key)
	//if err != nil {
	//
	//}
	//passwordBytes, err := util.GetSecrets(clientSet, namespace, emailConfig.PasswordSecret.Name, emailConfig.PasswordSecret.Key)
	//if err != nil {
	//
	//}
	//emailClient.Client = gomail.NewPlainDialer(emailConfig.SmtpHost, emailConfig.SmtpPort, string(userNameBytes), string(passwordBytes))
	////emailClient.Client.SSL = true
	//if emailClient.Client == nil {
	//	panic("Failed to create Email client")
	//
	//}
	return &emailClient
}

func (e *EmailClient) SendSuccessNotification(msg ...string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", e.Config.FromEmailId)
	m.SetHeader("To", e.Config.SenderList...)
	m.SetHeader("Subject", msg[1])
	m.SetBody("text/html", msg[0])
	// Send the email to Bob, Cora and Dan.
	if err := e.Client.DialAndSend(m); err != nil {
		panic(err)
	}

	return nil
}

func (e *EmailClient) SendWarningNotification(msg ...string) error {
	return nil
}
func (e *EmailClient) SendFailledNotification(msg ...string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", e.Config.FromEmailId)
	m.SetHeader("To", e.Config.SenderList...)
	m.SetHeader("Subject", msg[1])
	m.SetBody("text/html", msg[0])
	// Send the email to Bob, Cora and Dan.
	if err := e.Client.DialAndSend(m); err != nil {
		panic(err)
	}

	return nil
}
func (e *EmailClient) SendInfoNotification(msg ...string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", e.Config.FromEmailId)
	m.SetHeader("To", e.Config.SenderList...)
	m.SetHeader("Subject", msg[1])
	m.SetBody("text/html", msg[0])
	// Send the email to Bob, Cora and Dan.
	if err := e.Client.DialAndSend(m); err != nil {
		panic(err)
	}

	return nil
}
