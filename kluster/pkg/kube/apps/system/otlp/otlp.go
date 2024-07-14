package otlp

import (
	"context"
	"embed"
	"io/fs"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	log "github.com/sirupsen/logrus"
	"sigs.k8s.io/yaml"

	"github.com/ultary/monokube/kluster/pkg/helm"
	"github.com/ultary/monokube/kluster/pkg/k8s"
)

//go:embed helm/*
var chartFS embed.FS

func Enable(ctx context.Context, cluster *k8s.Cluster) {
	rootFS, _ := fs.Sub(chartFS, "helm")
	values := map[string]interface{}{}
	releaseName := "otlp"
	namespace := "kube-system"
	manifests := helm.BuildFromFileSystem(rootFS, values, releaseName, namespace)
	for _, m := range manifests {

		var raw map[string]interface{}
		if err := yaml.Unmarshal(m, &raw); err != nil {
			log.Fatalf("Failed to unmarshal YAML: %v", err)
		}

		obj := &unstructured.Unstructured{Object: raw}
		if err := cluster.ApplyUnstructured(ctx, obj, namespace); err != nil {
			log.Fatalf("%v", err)
		}
	}
}

func Disable(cluster *k8s.Cluster) {

}

func Update(cluster *k8s.Cluster) {

}
