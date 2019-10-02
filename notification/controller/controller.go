package controller

import (
	"bytes"
	"fmt"
	"github.com/antchfx/jsonquery"
	"github.com/argoproj-labs/argo-kube-notifier/notification/integration"
	"github.com/argoproj-labs/argo-kube-notifier/pkg/apis/argoproj/v1alpha1"
	"github.com/argoproj-labs/argo-kube-notifier/util"
	log "github.com/sirupsen/logrus"
	"html/template"
	apiv1 "k8s.io/api/core/v1"
	apierr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/json"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/util/workqueue"
	config2 "sigs.k8s.io/controller-runtime/pkg/client/config"
	"strings"
	"time"
)

const DELAY_WATCH_EVENT = "DELAY_WATCH"

type NewNotificationController struct {
	ObjectQueue        workqueue.RateLimitingInterface
	DelayQueue         workqueue.RateLimitingInterface
	ResourceMap        map[string]map[string]v1alpha1.Notification
	ResourceChanMap    map[string]chan struct{}
	NotifierMap        map[string]map[string]integration.NotifierInterface
	Namespace          string
	ResourceVersionMap map[string]string
}

func CreateNewNotificationController(namespace string) NewNotificationController {

	nnc := NewNotificationController{}
	nnc.ResourceMap = make(map[string]map[string]v1alpha1.Notification)
	nnc.NotifierMap = make(map[string]map[string]integration.NotifierInterface)
	nnc.ResourceChanMap = make(map[string]chan struct{})
	nnc.ObjectQueue = workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())
	nnc.ResourceVersionMap = make(map[string]string)
	nnc.Namespace = namespace
	return nnc
}

func (nnc *NewNotificationController) Starworker(chan1 chan struct{}) {
	newChannel := make(chan struct{})
	for i := 0; i < 5; i++ {
		go wait.Until(nnc.runWorker, time.Second, newChannel)
	}
	<-chan1
}

func (nnc *NewNotificationController) GenerateMapKey(resource schema.GroupVersionResource) string {
	return resource.Group + "_" + resource.Version + "_" + resource.Resource
}

func (nnc *NewNotificationController) UnRegisterNotification(notification *v1alpha1.Notification) {
	resourceName := nnc.GenerateMapKey(notification.Spec.MonitorResource)

	if notificationMap, ok := nnc.ResourceMap[resourceName]; ok {
		delete(notificationMap, notification.Name)
		//if()
	}
}

func (nnc *NewNotificationController) RegisterNotification(notification *v1alpha1.Notification) {

	resourceName := nnc.GenerateMapKey(notification.Spec.MonitorResource)

	if notificationMap, ok := nnc.ResourceMap[resourceName]; ok {
		notificationMap[notification.Name] = *notification
		nnc.ResourceMap[resourceName] = notificationMap
		nnc.NotifierMap[notification.Name] = nnc.getNotifier(*notification)
	} else {
		notificationMap = map[string]v1alpha1.Notification{}
		watcher := NewWatcher(&notification.Spec.MonitorResource, nnc.ObjectQueue, "", nnc.ResourceVersionMap)
		notificationMap[notification.Name] = *notification
		nnc.ResourceMap[resourceName] = notificationMap
		nnc.NotifierMap[notification.Name] = nnc.getNotifier(*notification)
		go watcher.watch()
	}
}

func (nnc *NewNotificationController) getNotifier(notification v1alpha1.Notification) map[string]integration.NotifierInterface {

	var notifierMap = make(map[string]integration.NotifierInterface)

	for _, notifier := range notification.Spec.Notifier {
		if notifier.Slack != nil {
			config, err := config2.GetConfig()
			if err != nil {

			}
			client, err := kubernetes.NewForConfig(config)
			namespace := "default"
			if notification.Namespace != "" {
				namespace = notification.Namespace
			}
			notifierMap[notifier.Name] = integration.NewSlackClient(client, namespace, notifier.Slack)
		}
		if notifier.Email != nil {
			config, err := config2.GetConfig()
			if err != nil {

			}
			client, err := kubernetes.NewForConfig(config)
			namespace := "default"
			if notification.Namespace != "" {
				namespace = notification.Namespace
			}
			notifierMap[notifier.Name] = integration.NewEmailClient(client, namespace, notifier.Email)
		}
	}
	return notifierMap
}

func (nnc *NewNotificationController) runWorker() {
	for nnc.processNextItem() {
	}
}

func (nnc *NewNotificationController) GenerateResourceNameKey(kind schema.GroupVersionKind) string {
	return kind.Group + "_" + kind.Version + "_" + kind.Kind
}
func (nnc *NewNotificationController) processNextItem() bool {
	object, quit := nnc.ObjectQueue.Get()
	if quit {
		return false
	}
	defer nnc.ObjectQueue.Done(object)
	event := object.(watch.Event)
	nnc.processObject(event, nnc.GenerateResourceNameKey(event.Object.GetObjectKind().GroupVersionKind()))

	return true
}

func (nnc *NewNotificationController) processObject(event watch.Event, resourceName string) {
	notifications, ok := nnc.ResourceMap[strings.ToLower(resourceName+"s")]
	if !ok {
		return
	}
	delayValidation := event.Type == DELAY_WATCH_EVENT
	for i := range notifications {
		notification := notifications[i]
		nnc.processingNotification(event, notification, delayValidation)
	}
	//for i := range notifications {
	//	notification := notifications[i]
	//	if notification.TriggerAction == "" {
	//		nnc.processingNotification(event, notification)
	//	} else {
	//
	//		switch notification.TriggerAction {
	//		case string(watch.Added):
	//			nnc.processingNotification(event, notification)
	//		case string(watch.Modified):
	//			nnc.processingNotification(event, notification)
	//		case string(watch.Deleted):
	//			nnc.processingNotification(event, notification)
	//		case string(watch.Error):
	//			nnc.processingNotification(event, notification)
	//		default:
	//			log.Info()
	//		}
	//	}
	//
	//}
}

func (nnc *NewNotificationController) processingNotification(event watch.Event, notification v1alpha1.Notification, delayValidation bool) {

	jsonObject, err := json.Marshal(event.Object)
	if err != nil {
		log.Error(err)
	}
	nnc.processRules(jsonObject, notification, delayValidation)

}

func (nnc *NewNotificationController) processRules(jsonByte []byte, notification v1alpha1.Notification, delayValidation bool) {
	doc, _ := jsonquery.Parse(strings.NewReader(string(jsonByte)))
	name := jsonquery.FindOne(doc, "metadata/name")
	for _, rule := range notification.Spec.Rules {

		if ValidateRule(&rule, doc) {
			if rule.InitialDelaySec > 0 && !delayValidation {
				go nnc.delayTrigger(rule.InitialDelaySec, notification.Spec.MonitorResource, notification.Namespace, name.InnerText())
				continue
			}
			log.Infof("Rule met condition. Event will be trigger. NotificationName=%s, Rule=%v", notification.Name, rule)

			nnc.processEvents(rule.Events, jsonByte, notification.Name)

		}

	}
}

func (nnc *NewNotificationController) delayTrigger(delay int, resource schema.GroupVersionResource, namespace, name string) {
	time.Sleep(time.Duration(delay) * time.Second)
	object, err := util.GetObject(resource, namespace, name)
	if err != nil {
		log.Warnf("Error occured getting resource. %v", err)
		return
	}
	event := watch.Event{
		Type:   DELAY_WATCH_EVENT,
		Object: object,
	}
	nnc.ObjectQueue.Add(event)

}

func (nnc *NewNotificationController) processEvents(events []v1alpha1.Event, jsonByte []byte, name string) {
	for _, event := range events {

		nnc.ProcessEvent(event, jsonByte, name)
	}
}

func (nnc *NewNotificationController) ProcessEvent(event v1alpha1.Event, jsonByte []byte, name string) {
	message := event.Message
	subject := event.EmailSubject
	if strings.Contains(message, "{{") || strings.Contains(subject, "{{") {
		jsonMap := map[string]interface{}{}
		if err := json.Unmarshal(jsonByte, &jsonMap); err != nil {
			panic(err)
		}
		message = nnc.SubsutiteString(event.Message, jsonMap)
		if event.EmailSubject != "" {
			subject = nnc.SubsutiteString(event.EmailSubject, jsonMap)
		}
	}
	notifierMap := nnc.NotifierMap[name]

	for _, notifierName := range event.NotifierNames {
		log.Debugf("Sending a message to %s", notifierName)
		if notifier, ok := notifierMap[notifierName]; ok {
			nnc.SendMessage(event.NotificationLevel, notifier, message, subject)
		} else {
			log.Warnf("Notifier not found: %s", notifierName)
		}
	}
}

func (nnc *NewNotificationController) SubsutiteString(rawmsg string, jsonMap map[string]interface{}) string {
	buff := bytes.NewBufferString("")
	t := template.New("")
	t.Parse(rawmsg)
	t.Execute(buff, jsonMap)
	return buff.String()
}

func (nc *NewNotificationController) SendMessage(notificationLevel string, integration integration.NotifierInterface, message ...string) {
	switch notificationLevel {
	case v1alpha1.NOTIFICATION_LEVEL_INFO:
		integration.SendSuccessNotification(message...)
	case v1alpha1.NOTIFICATION_LEVEL_WARNING:
		integration.SendWarningNotification(message...)
	case v1alpha1.NOTIFICATION_LEVEL_CRITICAL:
		integration.SendFailledNotification(message...)
	default:
		log.Warnf("Unknown notification level: %s", notificationLevel)
	}
}

func (nm *NewNotificationController) retrieveLastSyncResourceVersion() {
	config, err := config2.GetConfig() //clientcmd.BuildConfigFromFlags("", "")
	if err != nil {
		panic(err)
	}
	kubeclientset := kubernetes.NewForConfigOrDie(config)
	configMap, err := kubeclientset.CoreV1().ConfigMaps("default").Get("argo-argo-kube-notifier-resource-map", metav1.GetOptions{})

	nm.ResourceVersionMap = configMap.Data
}

func (nm *NewNotificationController) saveResourceVersion() {
	config, err := config2.GetConfig() //clientcmd.BuildConfigFromFlags("", "")
	if err != nil {
		panic(err.Error())
	}
	configMap := apiv1.ConfigMap{
		Data:       nm.ResourceVersionMap,
		ObjectMeta: metav1.ObjectMeta{Name: "argo-argo-kube-notifier-resource-map"},
	}
	kubeclientset := kubernetes.NewForConfigOrDie(config)
	_, err = kubeclientset.CoreV1().ConfigMaps("default").Create(&configMap)
	if apierr.IsAlreadyExists(err) {
		_, err = kubeclientset.CoreV1().ConfigMaps("default").Update(&configMap)
	}

	for k, v := range nm.ResourceVersionMap {
		fmt.Println(k, v)
	}
}
