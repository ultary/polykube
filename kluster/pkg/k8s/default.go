package k8s

import (
	"context"
	"encoding/json"
	stderrors "errors"
	"strings"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"

	"github.com/ultary/polykube/kluster/pkg/db/models"
)

////////////////////////////////////////////////////////////////
//
//  Apply (create or update) manifests
//

func (c *Cluster) ApplyUnstructured(ctx context.Context, obj *unstructured.Unstructured, namespace string) (err error) {

	discoveryClient := c.client.DiscoveryClient()

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

	dynamicClient := c.client.DynamicClient()

	var rc dynamic.ResourceInterface
	rc = dynamicClient.Resource(gvr)

	if isNamespaced {
		if n := obj.GetNamespace(); n != "" {
			namespace = n
		}
		rc = rc.(dynamic.NamespaceableResourceInterface).Namespace(namespace)
	}

	var appliedObj *unstructured.Unstructured
	opts := metav1.ApplyOptions{
		FieldManager: "None",
	}
	if appliedObj, err = rc.Apply(ctx, obj.GetName(), obj, opts); err != nil {
		log.Errorf("Failed to apply resource: %v", err)
		return err
	}

	metaObj, _ := meta.Accessor(appliedObj)
	metaObj.GetUID()
	metaObj.GetResourceVersion()

	log.Infoln("----------------")
	log.Infof("apiGroup: %s\n", appliedObj.GroupVersionKind().Group)
	log.Infof("apiVersion: %s\n", appliedObj.GetAPIVersion())
	log.Infof("kind: %s\n", appliedObj.GetKind())
	log.Infof("name: %s\n", metaObj.GetName())
	log.Infof("namespace: %s\n", metaObj.GetNamespace())
	log.Infof("uid: %s\n", metaObj.GetUID())
	log.Infof("resourceVersion: %s\n", metaObj.GetResourceVersion())

	uuid, err := uuid.Parse(string(appliedObj.GetUID()))
	if err != nil {
		log.Fatal(err)
	}

	row := &models.Resource{}
	err = c.db.Where("uid = ?", uuid).First(row).Error
	if err == nil {
		// TODO: update
		return nil
	}

	if stderrors.Is(err, gorm.ErrRecordNotFound) {
		var raw json.RawMessage
		raw, err = json.Marshal(obj)
		if err != nil {
			log.Fatal(err)
		}

		row = &models.Resource{
			APIGroup:   obj.GroupVersionKind().Group,
			APIVersion: obj.GetAPIVersion(),
			Kind:       obj.GetKind(),
			Name:       obj.GetName(),
			Namespace:  obj.GetNamespace(),
			Manifest:   raw,
			UID:        uuid,
		}
		if err = c.db.Create(row).Error; err != nil {
			log.Fatalf("Failed to create resource: %v", err)
		}

		return nil
	}

	log.Fatalf("Failed to find resource: %v", err)
	return nil
}

func (c *Cluster) ApplyNamespace(ctx context.Context, name string) (err error) {

	client := c.client.KubernetesClientset()

	namespace := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}

	result, err := client.CoreV1().Namespaces().Update(ctx, namespace, metav1.UpdateOptions{})
	if err == nil {
		log.Debug(result)
		return
	}
	if !errors.IsNotFound(err) {
		return
	}

	result, err = client.CoreV1().Namespaces().Create(ctx, namespace, metav1.CreateOptions{})
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

func (c *Cluster) ApplyDeployment(ctx context.Context, namespace string, deployment *appsv1.Deployment) (err error) {

	client := c.client.KubernetesClientset()

	if deployment.Namespace != "" {
		namespace = deployment.Namespace
	}

	var result *appsv1.Deployment
	result, err = client.AppsV1().Deployments(namespace).Update(ctx, deployment, metav1.UpdateOptions{})
	if err == nil {
		log.Debug(result)
		return
	}
	if !errors.IsNotFound(err) {
		return
	}

	result, err = client.AppsV1().Deployments(namespace).Create(ctx, deployment, metav1.CreateOptions{})
	if err == nil {
		log.Debug(result)
		return
	}
	if err != nil {
		log.Fatalf("Failed creating deployment: %v", err)
	}

	return
}

func (c *Cluster) ApplyStatefulSet(ctx context.Context, statefulSet *appsv1.StatefulSet, namespace string) (err error) {

	client := c.client.KubernetesClientset()

	if statefulSet.Namespace != "" {
		namespace = statefulSet.Namespace
	}

	var result *appsv1.StatefulSet
	result, err = client.AppsV1().StatefulSets(namespace).Update(ctx, statefulSet, metav1.UpdateOptions{})
	if err == nil {
		log.Debug(result)
		return
	}
	if !errors.IsNotFound(err) {
		return
	}

	result, err = client.AppsV1().StatefulSets(namespace).Create(ctx, statefulSet, metav1.CreateOptions{})
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

func (c *Cluster) ApplyConfigMap(ctx context.Context, namespace string, configmap *corev1.ConfigMap) (err error) {

	client := c.client.KubernetesClientset()

	if configmap.Namespace != "" {
		namespace = configmap.Namespace
	}

	var result *corev1.ConfigMap
	result, err = client.CoreV1().ConfigMaps(namespace).Update(ctx, configmap, metav1.UpdateOptions{})
	if err == nil {
		log.Debug(result)
		return
	}
	if !errors.IsNotFound(err) {
		return
	}

	result, err = client.CoreV1().ConfigMaps(namespace).Create(ctx, configmap, metav1.CreateOptions{})
	if err == nil {
		log.Debug(result)
		return
	}
	if err != nil {
		log.Fatalf("Failed createing configmap: %v", err)
	}

	return
}

func (c *Cluster) GetSecret(ctx context.Context, name, namespace string) (*corev1.Secret, error) {
	client := c.client.KubernetesClientset()
	return client.CoreV1().Secrets(namespace).Get(ctx, name, metav1.GetOptions{})
}

func (c *Cluster) CreateSecret(ctx context.Context, namespace string, secret *corev1.Secret) (*corev1.Secret, error) {
	if secret.Namespace != "" {
		namespace = secret.Namespace
	}
	client := c.client.KubernetesClientset()
	return client.CoreV1().Secrets(namespace).Create(ctx, secret, metav1.CreateOptions{})
}

// ---- Network ----

func (c *Cluster) ApplyService(ctx context.Context, service *corev1.Service, namespace string) (err error) {

	client := c.client.KubernetesClientset()

	if service.Namespace != "" {
		namespace = service.Namespace
	}

	var result *corev1.Service
	result, err = client.CoreV1().Services(namespace).Update(ctx, service, metav1.UpdateOptions{})
	if err == nil {
		log.Debug(result)
		return
	}
	if !errors.IsNotFound(err) {
		return
	}

	result, err = client.CoreV1().Services(namespace).Create(ctx, service, metav1.CreateOptions{})
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

func (c *Cluster) ApplyServiceAccount(ctx context.Context, sa *corev1.ServiceAccount, namespace string) (err error) {

	if sa.Namespace != "" {
		namespace = sa.Namespace
	}

	client := c.client.KubernetesClientset().CoreV1().ServiceAccounts(namespace)

	_, err = client.Get(ctx, sa.Name, metav1.GetOptions{})
	if err != nil {
		if !errors.IsNotFound(err) {
			log.Errorf("Failed to get ServiceAccount: %v", err)
			return
		}
		if _, err = client.Create(ctx, sa, metav1.CreateOptions{}); err != nil {
			log.Errorf("Failed to create ServiceAccount: %v", err)
			return
		}

	}

	// sa.ResourceVersion = current.ResourceVersion
	if _, err = client.Update(ctx, sa, metav1.UpdateOptions{}); err != nil {
		log.Errorf("Failed to update ServiceAccount: %v", err)
	}

	return
}

func (c *Cluster) ApplyClusterRole(ctx context.Context, cr *rbacv1.ClusterRole) (err error) {

	client := c.client.KubernetesClientset().RbacV1().ClusterRoles()

	_, err = client.Get(ctx, cr.Name, metav1.GetOptions{})
	if err != nil {
		if !errors.IsNotFound(err) {
			log.Errorf("Failed to get ClusterRole: %v", err)
			return
		}
		if _, err = client.Create(ctx, cr, metav1.CreateOptions{}); err != nil {
			log.Errorf("Failed to create ClusterRole: %v", err)
			return
		}

	}

	// cr.ResourceVersion = current.ResourceVersion
	if _, err = client.Update(ctx, cr, metav1.UpdateOptions{}); err != nil {
		log.Errorf("Failed to update ClusterRole: %v", err)
	}

	return
}

func (c *Cluster) ApplyClusterRoleBiding(ctx context.Context, crb *rbacv1.ClusterRoleBinding) (err error) {

	client := c.client.KubernetesClientset().RbacV1().ClusterRoleBindings()

	_, err = client.Get(ctx, crb.Name, metav1.GetOptions{})
	if err != nil {
		if !errors.IsNotFound(err) {
			log.Errorf("Failed to get ClusterRoleBinding: %v", err)
			return
		}
		if _, err = client.Create(ctx, crb, metav1.CreateOptions{}); err != nil {
			log.Errorf("Failed to create ClusterRoleBinding: %v", err)
			return
		}
	}

	// crb.ResourceVersion = current.ResourceVersion
	if _, err = client.Update(ctx, crb, metav1.UpdateOptions{}); err != nil {
		log.Errorf("Failed to update ClusterRoleBinding: %v", err)
	}

	return
}
