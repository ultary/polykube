package k8s

import (
	"github.com/ultary/polykube/kluster/pkg/k8s/ext/certmanager"
	"github.com/ultary/polykube/kluster/pkg/k8s/ext/istio"
)

func (c *Cluster) CertManager() *certmanager.Client {
	return certmanager.NewClient(c.client.certmanagerClientset)
}

func (c *Cluster) Istio() *istio.Client {
	return istio.NewClient(c.client.istioClientset)
}
