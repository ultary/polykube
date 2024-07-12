package monokube

import (
	log "github.com/sirupsen/logrus"

	"github.com/ultary/monokube/core/pkg/apps"
	"github.com/ultary/monokube/core/pkg/helm"
	"github.com/ultary/monokube/core/pkg/k8s"
)

func Install(ctx k8s.Context, chartPath, namespace string) {

	if err := k8s.ApplyNamespace(ctx, namespace); err != nil {
		log.Fatalf("error creating Namespace: %v", err)
	}

	resources := helm.Parse(chartPath, namespace)
	sequence := []apps.Resource{
		resources["gateway"],
		resources["kafka"],
		resources["minio"],
		resources["postgres"],
		//resources["otel_agent"],
	}

	for _, s := range sequence {
		s.Apply(ctx, namespace)
	}

	// for _, topic := range [...]string{"otlp_logs", "otlp_metrics", "otlp_spans"} {
	// 	if err := kafka.CreateTopic(ctx, client, topic); err != nil {
	// 		log.Fatalf("Kafka topic(%s) creation failed: %v", topic, err)
	// 	}
	// }
	// otlp.Apply()
	// otlp.ApplyConfigMap()
	// otlp.ApplyCollector()
	// postgres.Apply(ctx, client, manifests)
	// otlp.ApplyConsumer()
	// grafana.Apply(ctx, client, manifests)
}
