package otlp

import (
	"embed"
	"io/fs"

	log "github.com/sirupsen/logrus"
	"sigs.k8s.io/yaml"

	"github.com/ultary/monokube/kluster/pkg/helm"
	"github.com/ultary/monokube/kluster/pkg/k8s"
)

//go:embed helm/*
var chartFS embed.FS

func Enable(client *k8s.Client) {
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

		// obj := &unstructured.Unstructured{Object: raw}
		// if err := client.ApplyUnstructured(nil, obj, namespace); err != nil {
		// 	log.Fatalf("%v", err)
		// }
	}
}

func Disable(client *k8s.Client) {

}

func Update(client *k8s.Client) {

}
