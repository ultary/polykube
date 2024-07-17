package watch

import (
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/util/workqueue"

	"github.com/ultary/monokube/kluster/pkg/db/models"
	"github.com/ultary/monokube/kluster/pkg/k8s"
	"github.com/ultary/monokube/kluster/pkg/k8s/ext"
)

type startFunc func(<-chan struct{})
type shutdownFuc func()

type Tower struct {
	db        *gorm.DB
	informers []ext.Informer
	queue     workqueue.RateLimitingInterface
	stop      chan struct{}
}

func NewTower(cluster *k8s.Cluster, db *gorm.DB) *Tower {

	queue := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())

	informers := []ext.Informer{
		cluster.Informer(queue),
		cluster.CertManager().Informer(queue),
		cluster.Istio().Informer(queue),
	}

	retval := &Tower{
		db:        db,
		informers: informers,
		queue:     queue,
	}

	return retval
}

func (t *Tower) Watch() {

	t.stop = make(chan struct{})
	for _, informer := range t.informers {
		informer.Start(t.stop)
	}

	for {
		obj, shutdown := t.queue.Get()
		if shutdown {
			break
		}

		// key, err := cache.MetaNamespaceKeyFunc(obj)
		// if err != nil {
		// 	log.SetReportCaller(true)
		// 	log.Fatalf("Error")
		// }

		// switch obj.(type) {
		// case *core.ConfigMap:
		// 	log.Printf("Processing change to ConfigMap: %s\n", key)
		// case *core.Secret:
		// 	log.Printf("Processing change to Secret: %s\n", key)
		// case *apps.Deployment:
		// 	log.Printf("Processing change to Deployment: %s\n", key)
		// case *apps.StatefulSet:
		// 	log.Printf("Processing change to StatefulSet: %s\n", key)
		// case *apps.DaemonSet:
		// 	log.Printf("Processing change to DaemonSet: %s\n", key)
		// case *istio.Gateway:
		// 	log.Printf("Processing change to Istio/Gateway: %s\n", key)
		// case *istio.VirtualService:
		// 	log.Printf("Processing change to Istio/VirtualService: %s\n", key)
		// case *certmanager.Certificate:
		// 	log.Printf("Processing change to CertManager/Certificate: %s\n", key)
		// case *certmanager.ClusterIssuer:
		// 	log.Printf("Processing change to CertManager/ClusterIssuer: %s\n", key)
		// case *certmanager.Issuer:
		// 	log.Printf("Processing change to CertManager/Issuer: %s\n", key)
		// }

		updateLastResourceTypeChanged(t.db, obj)

		t.queue.Done(obj)
	}
}

func (t *Tower) Shutdown() {
	log.Info("[Tower] Stopping process")

	close(t.stop)
	for _, informer := range t.informers {
		informer.Shutdown()
	}

	log.Info("[Tower] Stopped process")
}

func updateLastResourceTypeChanged(db *gorm.DB, obj interface{}) error {
	metaObj, _ := meta.Accessor(obj)
	runtimeObj, _ := obj.(runtime.Object)
	typeObj, _ := meta.TypeAccessor(runtimeObj)

	log.Printf("ApiVersion: %s, Kind: %s, Name: %s, Namespace: %s, UID: %s, ResourceVersion: %s\n",
		typeObj.GetAPIVersion(),
		typeObj.GetKind(),
		metaObj.GetName(),
		metaObj.GetNamespace(),
		metaObj.GetUID(),
		metaObj.GetResourceVersion())

	rv, err := strconv.ParseUint(metaObj.GetResourceVersion(), 10, 64)
	if err != nil {
		log.Fatal(err)
	}

	row := &models.LatestRsourceKindVersion{
		ResourceVersion: rv,
		UpdatedAt:       time.Now(),
	}
	return db.Save(row).Error
}
