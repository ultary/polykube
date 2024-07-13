package kube

import "github.com/ultary/monokube/kluster/pkg/k8s"

type Cluster struct {
	platform *platform
	system   *system
}

func NewCluster(client *k8s.Client) *Cluster {
	return &Cluster{
		platform: NewPlatform(client),
		system:   NewSystem(client),
	}
}

func (c *Cluster) Platform() *platform {
	return c.platform
}

func (c *Cluster) System() *system {
	return c.system
}
