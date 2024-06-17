package monokube

import (
	"log"
	"ultary.co/kluster/pkg/apps"
	"ultary.co/kluster/pkg/k8s"
	"ultary.co/kluster/pkg/net"
	"ultary.co/kluster/pkg/repo"
)

func Install(ctx k8s.Context) {

	const namespace = "monokube"
	const chartPath = "/Users/ghilbut/work/workbench/ultary/monokube/.helm"

	if err := k8s.ApplyNamespace(ctx, namespace); err != nil {
		log.Fatalf("error creating Namespace: %v", err)
	}

	manifests := repo.NewManifests(chartPath)

	gateway := net.NewGateway(manifests)
	kafka := apps.NewKafka(manifests)
	minio := apps.NewMinIO(manifests)
	postgres := apps.NewPostgreSQL(manifests)

	resources := []k8s.Resource{
		&gateway,
		&kafka,
		&minio,
		&postgres,
	}

	for _, r := range resources {
		r.Apply(ctx, namespace)
	}

	// kafka.Apply(ctx, client, manifests)
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
