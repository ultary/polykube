package kube

import (
	"flag"
	"path/filepath"

	certmanagerversioned "github.com/cert-manager/cert-manager/pkg/client/clientset/versioned"
	log "github.com/sirupsen/logrus"
	istioversioned "istio.io/client-go/pkg/clientset/versioned"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type kube struct {
	config               *rest.Config
	kubernetesClientset  *kubernetes.Clientset
	dynamicClient        *dynamic.DynamicClient
	discoveryClient      *discovery.DiscoveryClient
	certmanagerClientset *certmanagerversioned.Clientset
	istioClientset       *istioversioned.Clientset
}

func NewKube() *kube {

	var err error

	retval := &kube{}

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

func (c *kube) Config() *rest.Config {
	return c.config
}

func (c *kube) KubernetesClientset() *kubernetes.Clientset {
	return c.kubernetesClientset
}

func (c *kube) DiscoveryClient() *discovery.DiscoveryClient {
	return c.discoveryClient
}

func (c *kube) DynamicClient() *dynamic.DynamicClient {
	return c.dynamicClient
}

func (c *kube) CertManagerClientset() *certmanagerversioned.Clientset {
	return c.certmanagerClientset
}

func (c *kube) IstioClientset() *istioversioned.Clientset {
	return c.istioClientset
}
