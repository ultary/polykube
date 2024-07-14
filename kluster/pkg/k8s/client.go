package k8s

import (
	"bytes"
	"context"
	"fmt"
	"net/http"

	certmanagerversioned "github.com/cert-manager/cert-manager/pkg/client/clientset/versioned"
	log "github.com/sirupsen/logrus"
	istioversioned "istio.io/client-go/pkg/clientset/versioned"
	core "k8s.io/api/core/v1"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
)

type Client struct {
	config               *rest.Config
	kubernetesClientset  *kubernetes.Clientset
	dynamicClient        *dynamic.DynamicClient
	discoveryClient      *discovery.DiscoveryClient
	certmanagerClientset *certmanagerversioned.Clientset
	istioClientset       *istioversioned.Clientset
}

func NewClient(kubeconfig, kubecontext string) *Client {

	var err error

	retval := &Client{}

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

func (c *Client) Config() *rest.Config {
	return c.config
}

func (c *Client) KubernetesClientset() *kubernetes.Clientset {
	return c.kubernetesClientset
}

func (c *Client) DiscoveryClient() *discovery.DiscoveryClient {
	return c.discoveryClient
}

func (c *Client) DynamicClient() *dynamic.DynamicClient {
	return c.dynamicClient
}

func (c *Client) CertManagerClientset() *certmanagerversioned.Clientset {
	return c.certmanagerClientset
}

func (c *Client) IstioClientset() *istioversioned.Clientset {
	return c.istioClientset
}

////////////////////////////////////////////////////////////////
//
//  Pod shell execution
//

func (c *Client) Exec(ctx context.Context, namespace, podName, containerName string, command []string) (string, string, error) {

	config := c.Config()
	client := c.KubernetesClientset()

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
		log.Errorf("error while creating SPDY executor: %v", err)
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
