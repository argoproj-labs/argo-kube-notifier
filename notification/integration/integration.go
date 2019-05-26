package integration

type NotifierInterface interface {
	SendSuccessNotification(msg ...string) error
	SendWarningNotification(msg ...string) error
	SendFailledNotification(msg ...string) error
	SendInfoNotification(msg ...string) error
}
