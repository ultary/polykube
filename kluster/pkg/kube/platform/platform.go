package platform

import "github.com/ultary/polykube/kluster/pkg/k8s"

type platform struct {
	client *k8s.Client
}

func NewPlatform(client *k8s.Client) *platform {
	return &platform{
		client: client,
	}
}
