package apps

import (
	"github.com/ultary/monokube/core/pkg/k8s"
	"helm.sh/helm/v3/pkg/chart"
)

type Resource interface {
	Apply(ctx k8s.Context, namespace string) error
}

type Manifests map[string]map[string][]byte

func NewManifests() Manifests {
	return make(Manifests)
}

func (m Manifests) Get(kind, name string) []byte {
	group := m[kind]
	if group == nil {
		return nil
	}
	return group[name]
}

func (m Manifests) Set(kind, name string, manifest []byte) {
	if m[kind] == nil {
		m[kind] = make(map[string][]byte)
	}
	m[kind][name] = manifest
}

type Dependencies map[string]*chart.Dependency

func NewDependencies() Dependencies {
	return make(Dependencies)
}

func (d Dependencies) Get(name string) *chart.Dependency {
	return d[name]
}

func (d Dependencies) Set(name string, c *chart.Dependency) {
	d[name] = c
}
