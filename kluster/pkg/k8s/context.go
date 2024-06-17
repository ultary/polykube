package k8s

import (
	"context"
	"flag"
	"path/filepath"
	"time"

	certmanagerversioned "github.com/cert-manager/cert-manager/pkg/client/clientset/versioned"
	log "github.com/sirupsen/logrus"
	istioversioned "istio.io/client-go/pkg/clientset/versioned"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
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
