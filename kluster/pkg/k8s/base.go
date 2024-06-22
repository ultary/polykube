package k8s

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	certmanagerversioned "github.com/cert-manager/cert-manager/pkg/client/clientset/versioned"
	log "github.com/sirupsen/logrus"
	istioversioned "istio.io/client-go/pkg/clientset/versioned"
	apps "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
	"k8s.io/client-go/util/homedir"
)

type Context interface {
	context.Context
	Config() *rest.Config
	KubernetesClientset() *kubernetes.Clientset
	DiscoveryClient() *discovery.DiscoveryClient
	DynamicClient() *dynamic.DynamicClient
	CertManagerClientset() *certmanagerversioned.Clientset
	IstioClientset() *istioversioned.Clientset
}

type contextImpl struct {
	base                 context.Context
	config               *rest.Config
	kubernetesClientset  *kubernetes.Clientset
	dynamicClient        *dynamic.DynamicClient
	discoveryClient      *discovery.DiscoveryClient
	certmanagerClientset *certmanagerversioned.Clientset
	istioClientset       *istioversioned.Clientset
}

func NewContext(ctx context.Context) Context {

	var err error

	retval := &contextImpl{
		base: ctx,
	}

	// --------

	kubeconfig, kubecontext := "", "k3s"
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = filepath.Join(home, ".kube", "config")
	}
	kubeconfig = *flag.String("kubeconfig", kubeconfig, "(optional) absolute path to the kubeconfig file")
	kubecontext = *flag.String("context", kubecontext, "The name of the kubeconfig context to use")
	flag.Parse()

	// kubeconfig 파일 로드
	rules := &clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfig}
	overrides := &clientcmd.ConfigOverrides{CurrentContext: kubecontext}
	configLoader := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		rules,
		overrides,
	)

	// REST config 생성
	if retval.config, err = configLoader.ClientConfig(); err != nil {
		log.Fatalf("Error creating Kubernetes client configuration: %v", err)
	}

	// Kubernetes 클라이언트 생성
	if retval.kubernetesClientset, err = kubernetes.NewForConfig(retval.config); err != nil {
		log.Fatalf("Error creating Kubernetes client: %v", err)
	}

	// Discovery 클라이언트 생성
	if retval.discoveryClient, err = discovery.NewDiscoveryClientForConfig(retval.config); err != nil {
		log.Fatalf("Error creating discovery client: %v", err)
	}

	// Dynamic 클라이언트 생성
	if retval.dynamicClient, err = dynamic.NewForConfig(retval.config); err != nil {
		log.Fatalf("Error creating dynamic client: %v", err)
	}

	// Cert-Manager 클라이언트 생성
	if retval.certmanagerClientset, err = certmanagerversioned.NewForConfig(retval.config); err != nil {
		log.Fatalf("Error creating Cert-Manager client: %v\n", err)
	}

	// Istio 클라이언트 생성
	if retval.istioClientset, err = istioversioned.NewForConfig(retval.config); err != nil {
		log.Fatalf("Error creating Istio client: %v\n", err)
	}

	return retval
}

// ---- Context interface methods ----

func (c *contextImpl) Deadline() (deadline time.Time, ok bool) {
	return c.base.Deadline()
}

func (c *contextImpl) Done() <-chan struct{} {
	return c.base.Done()
}

func (c *contextImpl) Err() error {
	return c.base.Err()
}

func (c *contextImpl) Value(key any) any {
	return c.base.Value(key)
}

// ---- Monokube context methods ----

func (c *contextImpl) Config() *rest.Config {
	return c.config
}

func (c *contextImpl) KubernetesClientset() *kubernetes.Clientset {
	return c.kubernetesClientset
}

func (c *contextImpl) DiscoveryClient() *discovery.DiscoveryClient {
	return c.discoveryClient
}

func (c *contextImpl) DynamicClient() *dynamic.DynamicClient {
	return c.dynamicClient
}

func (c *contextImpl) CertManagerClientset() *certmanagerversioned.Clientset {
	return c.certmanagerClientset
}

func (c *contextImpl) IstioClientset() *istioversioned.Clientset {
	return c.istioClientset
}

////////////////////////////////////////////////////////////////
//
//  Pod shell execution
//

func Exec(ctx Context, namespace, podName, containerName string, command []string) (string, string, error) {

	config := ctx.Config()
	client := ctx.KubernetesClientset()

	req := client.CoreV1().RESTClient().
		Post().
		Resource("pods").
		Namespace(namespace).
		Name(podName).
		SubResource("exec").
		VersionedParams(&core.PodExecOptions{
			Container: containerName,
			Command:   command,
			Stdin:     false,
			Stdout:    true,
			Stderr:    true,
			TTY:       false,
		}, scheme.ParameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(config, http.MethodPost, req.URL())
	if err != nil {
		fmt.Errorf("error while creating SPDY executor: %v", err)
		return "", "", err
	}

	var stdout, stderr bytes.Buffer
	err = exec.StreamWithContext(
		ctx,
		remotecommand.StreamOptions{
			Stdin:  nil,
			Stdout: &stdout,
			Stderr: &stderr,
			Tty:    false,
		})
	if err != nil {
		return "", "", fmt.Errorf("error in Stream: %v", err)
	}

	outs := stdout.String()
	errs := stderr.String()
	return outs, errs, nil
}

////////////////////////////////////////////////////////////////
//
//  Apply (create or update) manifests
//

func ApplyNamespace(ctx Context, name string) (err error) {

	client := ctx.KubernetesClientset()

	namespace := &core.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}

	result, err := client.CoreV1().Namespaces().Update(ctx, namespace, metav1.UpdateOptions{})
	if err == nil {
		log.Debug(result)
		return
	}
	status := err.(*errors.StatusError).ErrStatus
	if status.Code != http.StatusNotFound {
		return
	}

	result, err = client.CoreV1().Namespaces().Create(ctx, namespace, metav1.CreateOptions{})
	if err == nil {
		log.Debug(result)
		return
	}
	if err != nil {
		panic(err.(*errors.StatusError))
	}

	return
}

// ---- Workloads ----

func ApplyDeployment(ctx Context, namespace string, deployment *apps.Deployment) (err error) {

	client := ctx.KubernetesClientset()

	if deployment.Namespace != "" {
		namespace = deployment.Namespace
	}

	var result *apps.Deployment
	result, err = client.AppsV1().Deployments(namespace).Update(ctx, deployment, metav1.UpdateOptions{})
	if err == nil {
		log.Debug(result)
		return
	}
	e, ok := err.(*errors.StatusError)
	if !ok || e.Status().Code != http.StatusNotFound {
		return
	}

	result, err = client.AppsV1().Deployments(namespace).Create(ctx, deployment, metav1.CreateOptions{})
	if err == nil {
		log.Debug(result)
		return
	}
	if err != nil {
		log.Fatal(err.Error())
	}

	return
}

func ApplyStatefulSet(ctx Context, statefulSet *apps.StatefulSet, namespace string) (err error) {

	client := ctx.KubernetesClientset()

	if statefulSet.Namespace != "" {
		namespace = statefulSet.Namespace
	}

	var result *apps.StatefulSet
	result, err = client.AppsV1().StatefulSets(namespace).Update(ctx, statefulSet, metav1.UpdateOptions{})
	if err == nil {
		log.Debug(result)
		return
	}
	e, ok := err.(*errors.StatusError)
	if !ok || e.Status().Code != http.StatusNotFound {
		return
	}

	result, err = client.AppsV1().StatefulSets(namespace).Create(ctx, statefulSet, metav1.CreateOptions{})
	if err == nil {
		log.Debug(result)
		return
	}
	if err != nil {
		log.Fatal(err.Error())
	}

	return
}

// ---- Config ----

func ApplyConfigMap(ctx Context, namespace string, configmap *core.ConfigMap) (err error) {

	client := ctx.KubernetesClientset()

	if configmap.Namespace != "" {
		namespace = configmap.Namespace
	}

	var result *core.ConfigMap
	result, err = client.CoreV1().ConfigMaps(namespace).Update(ctx, configmap, metav1.UpdateOptions{})
	if err == nil {
		log.Debug(result)
		return
	}
	e, ok := err.(*errors.StatusError)
	if !ok || e.Status().Code != http.StatusNotFound {
		return
	}

	result, err = client.CoreV1().ConfigMaps(namespace).Create(ctx, configmap, metav1.CreateOptions{})
	if err == nil {
		log.Debug(result)
		return
	}
	if err != nil {
		log.Fatal(err.Error())
	}

	return
}

func GetSecret(ctx Context, name, namespace string) (*core.Secret, error) {
	client := ctx.KubernetesClientset()
	return client.CoreV1().Secrets(namespace).Get(ctx, name, metav1.GetOptions{})
}

func CreateSecret(ctx Context, namespace string, secret *core.Secret) (*core.Secret, error) {
	if secret.Namespace != "" {
		namespace = secret.Namespace
	}
	client := ctx.KubernetesClientset()
	return client.CoreV1().Secrets(namespace).Create(ctx, secret, metav1.CreateOptions{})
}

// ---- Network ----

func ApplyService(ctx Context, service *core.Service, namespace string) (err error) {

	client := ctx.KubernetesClientset()

	if service.Namespace != "" {
		namespace = service.Namespace
	}

	var result *core.Service
	result, err = client.CoreV1().Services(namespace).Update(ctx, service, metav1.UpdateOptions{})
	if err == nil {
		log.Debug(result)
		return
	}
	e, ok := err.(*errors.StatusError)
	if !ok || e.Status().Code != http.StatusNotFound {
		return
	}

	result, err = client.CoreV1().Services(namespace).Create(ctx, service, metav1.CreateOptions{})
	if err == nil {
		log.Debug(result)
		return
	}
	if err != nil {
		log.Fatal(err.Error())
	}

	return
}

// ---- Storage ----

// ---- Unstructured ----

func ApplyUnstructured(ctx Context, obj *unstructured.Unstructured, namespace string) (err error) {

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

	opts := metav1.ApplyOptions{
		FieldManager: "None",
	}
	if _, err = rc.Apply(ctx, obj.GetName(), obj, opts); err != nil {
		log.Errorf("Failed to apply resource: %v", err)
		return err
	}

	log.Infof("Successfully applied resource: %s/%s\n", gvr.Resource, obj.GetName())
	return nil
}
