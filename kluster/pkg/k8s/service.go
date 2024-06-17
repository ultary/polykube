package k8s

import (
	"net/http"

	log "github.com/sirupsen/logrus"
	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func ApplyService(ctx Context, service *core.Service, namespace string) (err error) {

	client := ctx.KubernetesClientset()

	if service.Namespace != "" {
		namespace = service.Namespace
	}

	var result *core.Service
	result, err = client.CoreV1().Services(namespace).Update(ctx, service, meta.UpdateOptions{})
	if err == nil {
		log.Debug(result)
		return
	}
	e, ok := err.(*errors.StatusError)
	if !ok || e.Status().Code != http.StatusNotFound {
		return
	}

	result, err = client.CoreV1().Services(namespace).Create(ctx, service, meta.CreateOptions{})
	if err == nil {
		log.Debug(result)
		return
	}
	if err != nil {
		log.Fatal(err.Error())
	}

	return
}
