package k8s

import (
	"net/http"

	log "github.com/sirupsen/logrus"
	apps "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

func ToDeploymentFromManifest(manifest []byte) (deployment *apps.Deployment, err error) {
	deployment = &apps.Deployment{}
	err = yaml.Unmarshal(manifest, deployment)
	return
}

func ApplyDeployment(ctx Context, namespace string, deployment *apps.Deployment) (err error) {

	client := ctx.KubernetesClientset()

	if deployment.Namespace != "" {
		namespace = deployment.Namespace
	}

	var result *apps.Deployment
	result, err = client.AppsV1().Deployments(namespace).Update(ctx, deployment, meta.UpdateOptions{})
	if err == nil {
		log.Debug(result)
		return
	}
	e, ok := err.(*errors.StatusError)
	if !ok || e.Status().Code != http.StatusNotFound {
		return
	}

	result, err = client.AppsV1().Deployments(namespace).Create(ctx, deployment, meta.CreateOptions{})
	if err == nil {
		log.Debug(result)
		return
	}
	if err != nil {
		log.Fatal(err.Error())
	}

	return
}
