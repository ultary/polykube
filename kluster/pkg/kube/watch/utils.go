package watch

import (
	certmanagerinformers "github.com/cert-manager/cert-manager/pkg/client/informers/externalversions"
	istioinformers "istio.io/client-go/pkg/informers/externalversions"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
)

type addEventHandlerFunc func(handler cache.ResourceEventHandler) (cache.ResourceEventHandlerRegistration, error)

func getAddEventHanders(factory informers.SharedInformerFactory) []addEventHandlerFunc {
	c := factory.Core().V1()
	a := factory.Apps().V1()
	return []addEventHandlerFunc{
		c.Secrets().Informer().AddEventHandler,
		c.ConfigMaps().Informer().AddEventHandler,
		a.Deployments().Informer().AddEventHandler,
		a.StatefulSets().Informer().AddEventHandler,
		a.DaemonSets().Informer().AddEventHandler,
	}
}

func getIstioAddEventHandlers(factory istioinformers.SharedInformerFactory) []addEventHandlerFunc {
	v1 := factory.Networking().V1()
	return []addEventHandlerFunc{
		v1.Gateways().Informer().AddEventHandler,
		v1.VirtualServices().Informer().AddEventHandler,
	}
}

func getCertManagerAddEventHandlers(factory certmanagerinformers.SharedInformerFactory) []addEventHandlerFunc {
	v1 := factory.Certmanager().V1()
	return []addEventHandlerFunc{
		v1.Certificates().Informer().AddEventHandler,
		v1.ClusterIssuers().Informer().AddEventHandler,
		v1.Issuers().Informer().AddEventHandler,
	}
}
