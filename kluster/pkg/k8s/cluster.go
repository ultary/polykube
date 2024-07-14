package k8s

import (
	"bytes"
	"context"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/remotecommand"
)

type Cluster struct {
	client *Client
	db     *gorm.DB
}

func NewCluster(client *Client, db *gorm.DB) *Cluster {

	return &Cluster{
		client: client,
		db:     db,
	}
}

////////////////////////////////////////////////////////////////
//
//  Pod shell execution
//

func (c *Cluster) Exec(ctx context.Context, namespace, podName, containerName string, command []string) (string, string, error) {

	config := c.client.Config()
	client := c.client.KubernetesClientset()

	req := client.CoreV1().RESTClient().
		Post().
		Resource("pods").
		Namespace(namespace).
		Name(podName).
		SubResource("exec").
		VersionedParams(&corev1.PodExecOptions{
			Container: containerName,
			Command:   command,
			Stdin:     false,
			Stdout:    true,
			Stderr:    true,
			TTY:       false,
		}, scheme.ParameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(config, http.MethodPost, req.URL())
	if err != nil {
		log.Errorf("error while creating SPDY executor: %v", err)
		return "", "", err
	}

	var stdout, stderr bytes.Buffer
	err = exec.StreamWithContext(
		ctx,
		remotecommand.StreamOptions{
			Stdin:  nil,
			Stdout: &stdout,
			Stderr: &stderr,
			Tty:    false,
		})
	if err != nil {
		return "", "", fmt.Errorf("error in Stream: %v", err)
	}

	outs := stdout.String()
	errs := stderr.String()
	return outs, errs, nil
}
