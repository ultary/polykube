package receivers

import (
	"fmt"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog/v2"
)

type secretReceiver struct {
	queue workqueue.RateLimitingInterface
}

func NewSecretReceiver(factory informers.SharedInformerFactory) Receiver {

	self := &secretReceiver{
		queue: workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter()),
	}

	informer := factory.Core().V1().Secrets()
	_, err := informer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    self.onAdd,
		UpdateFunc: self.onUpdate,
		DeleteFunc: self.onDelete,
	})
	if err != nil {
		klog.Fatalf("Failed event handler registration")
	}

	return self
}

func (r *secretReceiver) Run() {
	for {
		obj, shutdown := r.queue.Get()
		if shutdown {
			break
		}

		// Process the event
		key := obj.(string)
		fmt.Printf("Processing change to Secret: %s\n", key)
		r.queue.Done(obj)
	}
}

func (r *secretReceiver) onAdd(obj interface{}) {
	key, err := cache.MetaNamespaceKeyFunc(obj)
	if err == nil {
		r.queue.Add(key)
	}
}

func (r *secretReceiver) onUpdate(oldObj, newObj interface{}) {
	key, err := cache.MetaNamespaceKeyFunc(newObj)
	if err == nil {
		r.queue.Add(key)
	}
}

func (r *secretReceiver) onDelete(obj interface{}) {
	key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
	if err == nil {
		r.queue.Add(key)
	}
}
