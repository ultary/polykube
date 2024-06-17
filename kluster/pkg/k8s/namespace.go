package k8s

import (
	"net/http"

	log "github.com/sirupsen/logrus"
	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func ApplyNamespace(ctx Context, name string) (err error) {

	client := ctx.KubernetesClientset()

	namespace := &core.Namespace{
		ObjectMeta: meta.ObjectMeta{
			Name: name,
		},
	}

	result, err := client.CoreV1().Namespaces().Update(ctx, namespace, meta.UpdateOptions{})
	if err == nil {
		log.Debug(result)
		return
	}
	status := err.(*errors.StatusError).ErrStatus
	if status.Code != http.StatusNotFound {
		return
	}

	result, err = client.CoreV1().Namespaces().Create(ctx, namespace, meta.CreateOptions{})
	if err == nil {
		log.Debug(result)
		return
	}
	if err != nil {
		panic(err.(*errors.StatusError))
	}

	return
}
