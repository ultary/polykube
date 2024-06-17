package k8s

import (
	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetSecret(ctx Context, name, namespace string) (*core.Secret, error) {
	client := ctx.KubernetesClientset()
	return client.CoreV1().Secrets(namespace).Get(ctx, name, meta.GetOptions{})
}

func CreateSecret(ctx Context, namespace string, secret *core.Secret) (*core.Secret, error) {
	if secret.Namespace != "" {
		namespace = secret.Namespace
	}
	client := ctx.KubernetesClientset()
	return client.CoreV1().Secrets(namespace).Create(ctx, secret, meta.CreateOptions{})
}
