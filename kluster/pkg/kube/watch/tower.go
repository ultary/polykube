package watch

import (
	"slices"
	"time"

	certmanager "github.com/cert-manager/cert-manager/pkg/apis/certmanager/v1"
	certmanagerinformers "github.com/cert-manager/cert-manager/pkg/client/informers/externalversions"
	log "github.com/sirupsen/logrus"
	istio "istio.io/client-go/pkg/apis/networking/v1beta1"
	istioinformers "istio.io/client-go/pkg/informers/externalversions"
	apps "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"

	"github.com/ultary/monokube/kluster/pkg/k8s"
)

type startFunc func(<-chan struct{})
type shutdownFuc func()

type Tower struct {
	queue         workqueue.RateLimitingInterface
	stop          chan struct{}
	startFuncs    []startFunc
	shutdownFuncs []shutdownFuc
}

func NewTower(client *k8s.Client) *Tower {

	const resyncDuration = 30 * time.Second

	clientset := client.KubernetesClientset()
	factory := informers.NewSharedInformerFactory(clientset, resyncDuration)

	// cert-manager
	certmanagerClientset := client.CertManagerClientset()
	certmanagerFactory := certmanagerinformers.NewSharedInformerFactory(certmanagerClientset, resyncDuration)

	// istio
	istioClientset := client.IstioClientset()
	istioFactory := istioinformers.NewSharedInformerFactory(istioClientset, resyncDuration)

	retval := &Tower{
		queue: workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter()),
		startFuncs: []startFunc{
			factory.Start,
			certmanagerFactory.Start,
			istioFactory.Start,
		},
		shutdownFuncs: []shutdownFuc{
			istioFactory.Shutdown,
			certmanagerFactory.Shutdown,
			factory.Shutdown,
		},
	}

	adders := slices.Concat(
		getAddEventHanders(factory),
		getCertManagerAddEventHandlers(certmanagerFactory),
		getIstioAddEventHandlers(istioFactory))

	handler := cache.ResourceEventHandlerFuncs{
		AddFunc:    retval.onAdd,
		UpdateFunc: retval.onUpdate,
		DeleteFunc: retval.onDelete,
	}

	for _, AddEventHandler := range adders {
		if _, err := AddEventHandler(handler); err != nil {
			log.SetReportCaller(true)
			log.Fatalf("Failed event handler registration: %v", err)
		}
	}

	return retval
}

func (t *Tower) Watch() {

	t.stop = make(chan struct{})
	for _, start := range t.startFuncs {
		start(t.stop)
	}

	for {
		obj, shutdown := t.queue.Get()
		if shutdown {
			break
		}

		key, err := cache.MetaNamespaceKeyFunc(obj)
		if err != nil {
			log.SetReportCaller(true)
			log.Fatalf("Error")
		}

		switch obj.(type) {
		case *core.ConfigMap:
			log.Printf("Processing change to ConfigMap: %s\n", key)
		case *core.Secret:
			log.Printf("Processing change to Secret: %s\n", key)
		case *apps.Deployment:
			log.Printf("Processing change to Deployment: %s\n", key)
		case *apps.StatefulSet:
			log.Printf("Processing change to StatefulSet: %s\n", key)
		case *apps.DaemonSet:
			log.Printf("Processing change to DaemonSet: %s\n", key)
		case *istio.Gateway:
			log.Printf("Processing change to Istio/Gateway: %s\n", key)
		case *istio.VirtualService:
			log.Printf("Processing change to Istio/VirtualService: %s\n", key)
		case *certmanager.Certificate:
			log.Printf("Processing change to CertManager/Certificate: %s\n", key)
		case *certmanager.ClusterIssuer:
			log.Printf("Processing change to CertManager/ClusterIssuer: %s\n", key)
		case *certmanager.Issuer:
			log.Printf("Processing change to CertManager/Issuer: %s\n", key)
		}

		t.queue.Done(obj)
	}
}

func (t *Tower) Shutdown() {
	log.Info("[Tower] Stopping process")

	close(t.stop)
	for _, shutdown := range t.shutdownFuncs {
		shutdown()
	}

	log.Info("[Tower] Stopped process")
}

func (t *Tower) onAdd(obj interface{}) {
	// key, err := cache.MetaNamespaceKeyFunc(obj)
	// if err == nil {
	// 	//t.queue.Add(key)
	// 	log.Printf("Add: %s", key)
	// }
	t.queue.Add(obj)
}

func (t *Tower) onUpdate(oldObj, newObj interface{}) {
	// o, _ := cache.MetaNamespaceKeyFunc(oldObj)
	// n, _ := cache.MetaNamespaceKeyFunc(newObj)
	// log.Printf("Update: %s → %s", o, n)
	t.queue.Add(newObj)
}

func (t *Tower) onDelete(obj interface{}) {
	// key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
	// if err == nil {
	// 	//t.queue.Add(key)
	// 	log.Printf("Delete: %s", key)
	// }
	t.queue.Add(obj)
}
