package kube

import "github.com/ultary/monokube/kluster/pkg/k8s"

type platform struct {
	client *k8s.Client
}

func NewPlatform(client *k8s.Client) *platform {
	return &platform{
		client: client,
	}
}
