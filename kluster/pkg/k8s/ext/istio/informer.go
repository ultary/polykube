package istio

import (
	"gorm.io/gorm"
	"istio.io/client-go/pkg/clientset/versioned"
	"istio.io/client-go/pkg/informers/externalversions"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"

	"github.com/ultary/polykube/kluster/pkg/k8s/ext"
)

type Informer struct {
	clientset *versioned.Clientset
	db        *gorm.DB
	factory   externalversions.SharedInformerFactory
	queue     workqueue.RateLimitingInterface
}

func (c *Client) Informer(queue workqueue.RateLimitingInterface) ext.Informer {

	factory := externalversions.NewSharedInformerFactory(c.clientset, ext.DefaultResyncDuration)
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

	v1 := factory.Networking().V1()
	v1.Gateways().Informer().AddEventHandler(handler)
	v1.VirtualServices().Informer().AddEventHandler(handler)

	return &Informer{
		clientset: c.clientset,
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
