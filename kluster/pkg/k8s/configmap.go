package k8s

import (
	"net/http"

	log "github.com/sirupsen/logrus"
	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

func ToConfigMapFromManifest(manifest []byte) (configmap *core.ConfigMap, err error) {
	configmap = &core.ConfigMap{}
	err = yaml.Unmarshal(manifest, configmap)
	return
}

func ApplyConfigMap(ctx Context, namespace string, configmap *core.ConfigMap) (err error) {

	client := ctx.KubernetesClientset()

	if configmap.Namespace == "" {
		configmap.Namespace = namespace
	} else {
		namespace = configmap.Namespace
	}

	var result *core.ConfigMap
	result, err = client.CoreV1().ConfigMaps(namespace).Update(ctx, configmap, meta.UpdateOptions{})
	if err == nil {
		log.Debug(result)
		return
	}
	e, ok := err.(*errors.StatusError)
	if !ok || e.Status().Code != http.StatusNotFound {
		return
	}

	result, err = client.CoreV1().ConfigMaps(namespace).Create(ctx, configmap, meta.CreateOptions{})
	if err == nil {
		log.Debug(result)
		return
	}
	if err != nil {
		log.Fatal(err.Error())
	}

	return
}
