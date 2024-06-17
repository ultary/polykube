package k8s

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	certmanagerversioned "github.com/cert-manager/cert-manager/pkg/client/clientset/versioned"
	log "github.com/sirupsen/logrus"
	istioversioned "istio.io/client-go/pkg/clientset/versioned"
	apps "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
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
	CertManagerClientset() *certmanagerversioned.Clientset
	IstioClientset() *istioversioned.Clientset
}

type contextImpl struct {
	base   context.Context
	config *rest.Config
	k      *kubernetes.Clientset
	c      *certmanagerversioned.Clientset
	i      *istioversioned.Clientset
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
	if retval.k, err = kubernetes.NewForConfig(retval.config); err != nil {
		log.Fatalf("Error creating Kubernetes client: %v", err)
	}

	// Cert-Manager 클라이언트 생성
	if retval.c, err = certmanagerversioned.NewForConfig(retval.config); err != nil {
		log.Fatalf("Error creating Cert-Manager client: %v\n", err)
	}

	// Istio 클라이언트 생성
	if retval.i, err = istioversioned.NewForConfig(retval.config); err != nil {
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
	return c.k
}

func (c *contextImpl) CertManagerClientset() *certmanagerversioned.Clientset {
	return c.c
}

func (c *contextImpl) IstioClientset() *istioversioned.Clientset {
	return c.i
}

////////////////////////////////////////////////////////////////
//
//
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
//  Manifests API
//

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

// ---- Workloads ----

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

// ---- Config ----

func ApplyConfigMap(ctx Context, namespace string, configmap *core.ConfigMap) (err error) {

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

// ---- Network ----

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

// ---- Storage ----
