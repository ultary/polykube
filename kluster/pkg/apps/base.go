package apps

import "ultary.co/kluster/pkg/k8s"

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
