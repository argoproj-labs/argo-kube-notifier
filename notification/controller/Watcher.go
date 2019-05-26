package controller

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/util/workqueue"
	config2 "sigs.k8s.io/controller-runtime/pkg/client/config"
)

type Watcher struct {
	resource        *schema.GroupVersionResource
	objectQueue     workqueue.RateLimitingInterface
	apiWatcher      watch.Interface
	lastSyncVersion string
	resourceVerMap  map[string]string
}

func NewWatcher(resource *schema.GroupVersionResource, objectQueue workqueue.RateLimitingInterface, lastSyncVersion string, resourceVerMap map[string]string) Watcher {
	nw := Watcher{
		resource:        resource,
		objectQueue:     objectQueue,
		resourceVerMap:  resourceVerMap,
		lastSyncVersion: lastSyncVersion,
	}
	return nw
}

func (w *Watcher) watch() {
	if w.resource.Resource == "" {
		return
	}
	for {
		w.createWatcher()
		w.runWatch()
	}
}

func (w *Watcher) createWatcher() {
	config, err := config2.GetConfig()
	if err != nil {
		panic(err.Error())
	}
	clientset, err := dynamic.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	api := clientset.Resource(*w.resource)

	listStruct, err := api.List(v1.ListOptions{})
	if err != nil || listStruct == nil {
		return
	}
	w.lastSyncVersion = listStruct.GetResourceVersion()
	fmt.Println(w.lastSyncVersion)

	w.apiWatcher, err = api.
		Watch(v1.ListOptions{ResourceVersion: w.lastSyncVersion})
	if err != nil {
		log.Fatal(err)
	}
}

func (w *Watcher) runWatch() {
	ch := w.apiWatcher.ResultChan()
	for event := range ch {
		w.objectQueue.Add(event)

	}
}
