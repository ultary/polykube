package k8s

import (
	"net/http"

	log "github.com/sirupsen/logrus"
	apps "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func ApplyStatefulSet(ctx Context, statefulSet *apps.StatefulSet, namespace string) (err error) {

	client := ctx.KubernetesClientset()

	if statefulSet.Namespace != "" {
		namespace = statefulSet.Namespace
	}

	var result *apps.StatefulSet
	result, err = client.AppsV1().StatefulSets(namespace).Update(ctx, statefulSet, meta.UpdateOptions{})
	if err == nil {
		log.Debug(result)
		return
	}
	e, ok := err.(*errors.StatusError)
	if !ok || e.Status().Code != http.StatusNotFound {
		return
	}

	result, err = client.AppsV1().StatefulSets(namespace).Create(ctx, statefulSet, meta.CreateOptions{})
	if err == nil {
		log.Debug(result)
		return
	}
	if err != nil {
		log.Fatal(err.Error())
	}

	return
}
