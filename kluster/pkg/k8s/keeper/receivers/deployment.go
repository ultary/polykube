package receivers

import (
	"fmt"

	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog/v2"
)

type deploymentReceiver struct {
	queue workqueue.RateLimitingInterface
}

func NewDeploymentReceiver(factory informers.SharedInformerFactory) Receiver {

	self := &deploymentReceiver{
		queue: workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter()),
	}

	informer := factory.Apps().V1().Deployments()
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

func (r *deploymentReceiver) Run() {
	for {
		obj, shutdown := r.queue.Get()
		if shutdown {
			break
		}

		// Process the event
		key := obj.(string)
		fmt.Printf("Processing change to Deployment: %s\n", key)
		r.queue.Done(obj)
	}
}

func (r *deploymentReceiver) onAdd(obj interface{}) {
	key, err := cache.MetaNamespaceKeyFunc(obj)
	if err == nil {
		r.queue.Add(key)
	}

	// deployment := obj.(*v1.Deployment)
	// containers := deployment.Spec.Template.Spec.Containers
	// for _, container := range containers {
	// 	for _, env := range container.Env {
	// 		from := env.ValueFrom
	// 		cm := from.ConfigMapKeyRef
	// 		if cm != nil {
	// 			klog.Infof("Deployment - %s/%s - has configmap named %s", deployment.Namespace, deployment.Name, cm.Name)
	// 		}
	// 		s := from.SecretKeyRef
	// 		if s != nil {
	// 			klog.Infof("Deployment - %s/%s - has secret named %s", deployment.Namespace, deployment.Name, s.Name)
	// 		}
	// 	}
	// }
}

func (r *deploymentReceiver) onUpdate(oldObj, newObj interface{}) {
	key, err := cache.MetaNamespaceKeyFunc(newObj)
	if err == nil {
		r.queue.Add(key)
	}
}

func (r *deploymentReceiver) onDelete(obj interface{}) {
	key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
	if err == nil {
		r.queue.Add(key)
	}
}
