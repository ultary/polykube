package kafka

import (
	"context"
	"errors"

	log "github.com/sirupsen/logrus"

	"ultary.co/kluster/pkg/k8s"
)

func CreateTopic(ctx context.Context, client *k8s.Client, topic string) error {

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

	stdout, stderr, err := client.Exec(ctx, namespace, podName, containerName, command)
	if err != nil {
		log.Fatal(err)
	}
	if stderr != "" {
		return errors.New(stderr)
	}

	log.Info(stdout)
	return nil
}
