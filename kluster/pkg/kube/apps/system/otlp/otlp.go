package otlp

import (
	_ "embed"
	"github.com/ultary/monokube/kluster/pkg/k8s"
)

//go:embed manifests/sa.yaml
var saYaml []byte

//go:embed manifests/cr.yaml
var crYaml []byte

//go:embed manifests/crb.yaml
var crbYaml []byte

//go:embed manifests/otlp-agent.yaml
var agentValues []byte

//go:embed manifests/otlp-collector.yaml
var collectorValues []byte

func Enable(client *k8s.Client) {

}

func Disable(client *k8s.Client) {

}

func Update(client *k8s.Client) {

}
