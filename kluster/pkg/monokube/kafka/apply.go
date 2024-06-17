package kafka

import (
	"context"
	"log"

	"ultary.co/kluster/pkg/k8s"
	"ultary.co/kluster/pkg/repo"
)

func Apply(ctx context.Context, client *k8s.Client, manifests *repo.Manifests) error {

	const namespace = "monokube"

	// ==== StatefulSet ====

	raw := manifests.Get("StatefulSet", "kafka")
	if sts, err := k8s.ToStatefulSetFromManifest(raw); err == nil {
		client.ApplyStatefulSet(ctx, namespace, sts)
	} else {
		log.Fatalf("Error unmarshalling YAML to StatefulSet: %v", err)
	}

	// ==== Service ====

	raw = manifests.Get("Service", "kafka")
	if service, err := k8s.ToServiceFromManifest(raw); err == nil {
		client.ApplyService(ctx, namespace, service)
	} else {
		log.Fatalf("Error unmarshalling YAML to StatefulSet: %v", err)
	}

	return nil
}
