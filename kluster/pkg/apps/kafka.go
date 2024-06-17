package apps

import (
	log "github.com/sirupsen/logrus"
	istio "istio.io/client-go/pkg/apis/networking/v1beta1"
	apps "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	"sigs.k8s.io/yaml"

	"ultary.co/kluster/pkg/k8s"
	"ultary.co/kluster/pkg/repo"
)

type Kafka struct {
	sts apps.StatefulSet
	sv  core.Service
	vsv istio.VirtualService
}

func NewKafka(manifests *repo.Manifests) (retval Kafka) {

	const name = "kafka"

	m := manifests.Get("StatefulSet", name)
	if err := yaml.Unmarshal(m, &retval.sts); err != nil {
		log.Fatalf("Error unmarshalling YAML to StatefulSet: %v", err)
	}

	m = manifests.Get("Service", name)
	if err := yaml.Unmarshal(m, &retval.sv); err != nil {
		log.Fatalf("Error unmarshalling YAML to Service: %v", err)
	}

	m = manifests.Get("VirtualService", name)
	if err := yaml.Unmarshal(m, &retval.vsv); err != nil {
		log.Fatalf("Error unmarshalling YAML to VirtualService: %v", err)
	}

	return
}

func (k *Kafka) Apply(ctx k8s.Context, namespace string) error {

	if err := k8s.ApplyStatefulSet(ctx, &k.sts, namespace); err != nil {
		log.Fatalf("Error applying Kafka StatefulSet: %v", err)
	}

	if err := k8s.ApplyService(ctx, &k.sv, namespace); err != nil {
		log.Fatalf("Error applying Kafka Service: %v", err)
	}

	if err := k8s.ApplyVirtualService(ctx, &k.vsv, namespace); err != nil {
		log.Fatalf("Error applying Kafka VirtualService: %v", err)
	}

	return nil
}
