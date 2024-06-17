package monokube

import (
	log "github.com/sirupsen/logrus"

	"ultary.co/kluster/pkg/apps"
	"ultary.co/kluster/pkg/apps/net"
	"ultary.co/kluster/pkg/helm"
	"ultary.co/kluster/pkg/k8s"
)

func Install(ctx k8s.Context, chartPath, namespace string) {

	if err := k8s.ApplyNamespace(ctx, namespace); err != nil {
		log.Fatalf("error creating Namespace: %v", err)
	}

	chart := helm.NewChart(chartPath, namespace)

	gateway := net.NewGateway(chart)
	kafka := apps.NewKafka(chart)
	minio := apps.NewMinIO(chart)
	postgres := apps.NewPostgreSQL(chart)

	resources := []apps.Resource{
		&gateway,
		&kafka,
		&minio,
		&postgres,
	}

	for _, r := range resources {
		r.Apply(ctx, namespace)
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
