package k8s

import (
	"strings"

	log "github.com/sirupsen/logrus"
	apps "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	rbac "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

////////////////////////////////////////////////////////////////
//
//  Apply (create or update) manifests
//

func (c *Client) ApplyUnstructured(ctx Context, obj *unstructured.Unstructured, namespace string) (err error) {

	discoveryClient := ctx.DiscoveryClient()

	gvk := obj.GroupVersionKind()

	gv := gvk.GroupVersion().String()
	resources, err := discoveryClient.ServerResourcesForGroupVersion(gv)
	if err != nil {
		log.Fatalf("Failed to get API resources: %v", err)
	}

	var gvr schema.GroupVersionResource
	var isNamespaced bool
	for _, rc := range resources.APIResources {
		const suffix = "/status"
		if rc.Kind == gvk.Kind && !strings.HasSuffix(rc.Name, suffix) {
			gvr.Group = gvk.Group
			gvr.Version = gvk.Version
			gvr.Resource = rc.Name
			isNamespaced = rc.Namespaced
		}
	}

	dynamicClient := ctx.DynamicClient()

	var rc dynamic.ResourceInterface
	rc = dynamicClient.Resource(gvr)

	if isNamespaced {
		if n := obj.GetNamespace(); n != "" {
			namespace = n
		}
		rc = rc.(dynamic.NamespaceableResourceInterface).Namespace(namespace)
	}

	opts := meta.ApplyOptions{
		FieldManager: "None",
	}
	if _, err = rc.Apply(ctx, obj.GetName(), obj, opts); err != nil {
		log.Errorf("Failed to apply resource: %v", err)
		return err
	}

	log.Infof("Successfully applied resource: %s/%s\n", gvr.Resource, obj.GetName())
	return nil
}

func (c *Client) ApplyNamespace(ctx Context, name string) (err error) {

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
	if !errors.IsNotFound(err) {
		return
	}

	result, err = client.CoreV1().Namespaces().Create(ctx, namespace, meta.CreateOptions{})
	if err == nil {
		log.Debug(result)
		return
	}
	if err != nil {
		log.Fatalf("Failed creating namespace: %v", err)
	}

	return
}

// ---- Workloads ----

func (c *Client) ApplyDeployment(ctx Context, namespace string, deployment *apps.Deployment) (err error) {

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
	if !errors.IsNotFound(err) {
		return
	}

	result, err = client.AppsV1().Deployments(namespace).Create(ctx, deployment, meta.CreateOptions{})
	if err == nil {
		log.Debug(result)
		return
	}
	if err != nil {
		log.Fatalf("Failed creating deployment: %v", err)
	}

	return
}

func (c *Client) ApplyStatefulSet(ctx Context, statefulSet *apps.StatefulSet, namespace string) (err error) {

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
	if !errors.IsNotFound(err) {
		return
	}

	result, err = client.AppsV1().StatefulSets(namespace).Create(ctx, statefulSet, meta.CreateOptions{})
	if err == nil {
		log.Debug(result)
		return
	}
	if err != nil {
		log.Fatalf("Failed creating statefulset: %v", err)
	}

	return
}

// ---- Config ----

func (c *Client) ApplyConfigMap(ctx Context, namespace string, configmap *core.ConfigMap) (err error) {

	client := ctx.KubernetesClientset()

	if configmap.Namespace != "" {
		namespace = configmap.Namespace
	}

	var result *core.ConfigMap
	result, err = client.CoreV1().ConfigMaps(namespace).Update(ctx, configmap, meta.UpdateOptions{})
	if err == nil {
		log.Debug(result)
		return
	}
	if !errors.IsNotFound(err) {
		return
	}

	result, err = client.CoreV1().ConfigMaps(namespace).Create(ctx, configmap, meta.CreateOptions{})
	if err == nil {
		log.Debug(result)
		return
	}
	if err != nil {
		log.Fatalf("Failed createing configmap: %v", err)
	}

	return
}

func (c *Client) GetSecret(ctx Context, name, namespace string) (*core.Secret, error) {
	client := ctx.KubernetesClientset()
	return client.CoreV1().Secrets(namespace).Get(ctx, name, meta.GetOptions{})
}

func (c *Client) CreateSecret(ctx Context, namespace string, secret *core.Secret) (*core.Secret, error) {
	if secret.Namespace != "" {
		namespace = secret.Namespace
	}
	client := ctx.KubernetesClientset()
	return client.CoreV1().Secrets(namespace).Create(ctx, secret, meta.CreateOptions{})
}

// ---- Network ----

func (c *Client) ApplyService(ctx Context, service *core.Service, namespace string) (err error) {

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
	if !errors.IsNotFound(err) {
		return
	}

	result, err = client.CoreV1().Services(namespace).Create(ctx, service, meta.CreateOptions{})
	if err == nil {
		log.Debug(result)
		return
	}
	if err != nil {
		log.Fatalf("Failed createing service: %v", err)
	}

	return
}

// ---- Storage ----

// ---- Access Control ----

func (c *Client) ApplyServiceAccount(ctx Context, sa *core.ServiceAccount, namespace string) (err error) {

	if sa.Namespace != "" {
		namespace = sa.Namespace
	}

	client := ctx.KubernetesClientset().CoreV1().ServiceAccounts(namespace)

	_, err = client.Get(ctx, sa.Name, meta.GetOptions{})
	if err != nil {
		if !errors.IsNotFound(err) {
			log.Errorf("Failed to get ServiceAccount: %v", err)
			return
		}
		if _, err = client.Create(ctx, sa, meta.CreateOptions{}); err != nil {
			log.Errorf("Failed to create ServiceAccount: %v", err)
			return
		}

	}

	// sa.ResourceVersion = current.ResourceVersion
	if _, err = client.Update(ctx, sa, meta.UpdateOptions{}); err != nil {
		log.Errorf("Failed to update ServiceAccount: %v", err)
	}

	return
}

func (c *Client) ApplyClusterRole(ctx Context, cr *rbac.ClusterRole) (err error) {

	client := ctx.KubernetesClientset().RbacV1().ClusterRoles()

	_, err = client.Get(ctx, cr.Name, meta.GetOptions{})
	if err != nil {
		if !errors.IsNotFound(err) {
			log.Errorf("Failed to get ClusterRole: %v", err)
			return
		}
		if _, err = client.Create(ctx, cr, meta.CreateOptions{}); err != nil {
			log.Errorf("Failed to create ClusterRole: %v", err)
			return
		}

	}

	// cr.ResourceVersion = current.ResourceVersion
	if _, err = client.Update(ctx, cr, meta.UpdateOptions{}); err != nil {
		log.Errorf("Failed to update ClusterRole: %v", err)
	}

	return
}

func (c *Client) ApplyClusterRoleBiding(ctx Context, crb *rbac.ClusterRoleBinding) (err error) {

	client := ctx.KubernetesClientset().RbacV1().ClusterRoleBindings()

	_, err = client.Get(ctx, crb.Name, meta.GetOptions{})
	if err != nil {
		if !errors.IsNotFound(err) {
			log.Errorf("Failed to get ClusterRoleBinding: %v", err)
			return
		}
		if _, err = client.Create(ctx, crb, meta.CreateOptions{}); err != nil {
			log.Errorf("Failed to create ClusterRoleBinding: %v", err)
			return
		}
	}

	// crb.ResourceVersion = current.ResourceVersion
	if _, err = client.Update(ctx, crb, meta.UpdateOptions{}); err != nil {
		log.Errorf("Failed to update ClusterRoleBinding: %v", err)
	}

	return
}
