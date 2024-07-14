package apps

/*
import (
	"errors"

	log "github.com/sirupsen/logrus"
	istio "istio.io/client-go/pkg/apis/networking/v1"
	apps "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	"sigs.k8s.io/yaml"

	"github.com/ultary/monokube/kluster/pkg/k8s"
)

type Kafka struct {
	sts apps.StatefulSet
	sv  core.Service
	vsv istio.VirtualService
}

func NewKafka(manifests Manifests) *Kafka {

	const name = "kafka"

	var retval Kafka

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

	return &retval
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

func CreateTopic(ctx k8s.Context, topic string) error {

	const (
		namespace     = "monokube"
		podName       = "kafka-0"
		containerName = "kafka"
	)

	command := []string{
		"kafka-topics",
		"--bootstrap-server",
		"localhost:9092",
		"--create",
		"--if-not-exists",
		"--partitions",
		"2",
		"--replication-factor",
		"1",
		"--topic",
		topic,
	}

	stdout, stderr, err := k8s.Exec(ctx, namespace, podName, containerName, command)
	if err != nil {
		log.Fatal(err)
	}
	if stderr != "" {
		return errors.New(stderr)
	}

	log.Info(stdout)
	return nil
}
*/
