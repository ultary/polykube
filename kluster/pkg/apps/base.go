package apps

import (
	"ultary.co/kluster/pkg/k8s"
)

type Resource interface {
	Apply(ctx k8s.Context, namespace string) error
}
