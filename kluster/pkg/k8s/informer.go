package k8s

import (
	"gorm.io/gorm"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"

	"github.com/ultary/polykube/kluster/pkg/k8s/ext"
)

type Informer struct {
	clientset kubernetes.Interface
	db        *gorm.DB
	factory   informers.SharedInformerFactory
	queue     workqueue.RateLimitingInterface
}

func (c *Cluster) Informer(queue workqueue.RateLimitingInterface) ext.Informer {

	factory := informers.NewSharedInformerFactory(c.client.kubernetesClientset, ext.DefaultResyncDuration)
	handler := cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			queue.Add(obj)
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			queue.Add(newObj)
		},
		DeleteFunc: func(obj interface{}) {
			queue.Add(obj)
		},
	}

	core := factory.Core().V1()
	core.Secrets().Informer().AddEventHandler(handler)
	core.ConfigMaps().Informer().AddEventHandler(handler)
	apps := factory.Apps().V1()
	apps.Deployments().Informer().AddEventHandler(handler)
	apps.StatefulSets().Informer().AddEventHandler(handler)
	apps.DaemonSets().Informer().AddEventHandler(handler)

	return &Informer{
		clientset: c.client.kubernetesClientset,
		factory:   factory,
		queue:     queue,
	}
}

func (i *Informer) Start(stopCh <-chan struct{}) {
	i.factory.Start(stopCh)
	i.factory.WaitForCacheSync(stopCh)
}

func (i *Informer) Shutdown() {
	i.factory.Shutdown()
}
