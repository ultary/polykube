package ext

import (
	"time"
)

const DefaultResyncDuration = 30 * time.Second

type Informer interface {
	//AddEventHandler(handler cache.ResourceEventHandler)
	Start(stopCh <-chan struct{})
	Shutdown()
	//OnAdd(obj interface{}, _ bool)
	//OnUpdate(oldObj, newObj interface{})
	//OnDelete(obj interface{})
}
