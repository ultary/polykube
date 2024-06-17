package repo

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/engine"
	"sigs.k8s.io/yaml"
)

type Manifests struct {
	manifests    map[string]map[string][]byte
	dependencies map[string]*chart.Dependency
	values       map[string]interface{}
}

func NewManifests(chartPath string) *Manifests {

	releaseName := "monokube"
	namespace := "monokube"
	valuesFile := filepath.Join(chartPath, "values.yaml")

	// Load the chart
	chartRequested, err := loader.Load(chartPath)
	if err != nil {
		log.Fatalf("Error loading chart: %v\n", err)
	}

	// Load values file if provided
	vals := map[string]interface{}{}
	if valuesFile != "" {
		vals, err = readValuesFile(valuesFile)
		if err != nil {
			fmt.Printf("Error reading values file: %v\n", err)
			os.Exit(1)
		}
	}

	////////////////////////////////////////////////////////////////
	//  Render

	// Create a rendering environment
	valuesToRender, err := chartutil.ToRenderValues(chartRequested, vals, chartutil.ReleaseOptions{
		Name:      releaseName,
		Namespace: namespace,
	}, nil)
	if err != nil {
		fmt.Printf("Error creating render values: %v\n", err)
		os.Exit(1)
	}

	// Render the chart
	rendered, err := engine.Render(chartRequested, valuesToRender)
	if err != nil {
		fmt.Printf("Error rendering chart: %v\n", err)
		os.Exit(1)
	}

	manifests := &Manifests{
		manifests: make(map[string]map[string][]byte),
		values:    vals,
	}

	// Print the rendered templates
	for _, content := range rendered {
		blocks := SplitYAML([]byte(content))
		for _, block := range blocks {
			var m KubernetesManifest
			yaml.Unmarshal(block, &m)
			manifests.Set(m.Kind, m.Metadata.Name, block)
		}
	}

	////////////////////////////////////////////////////////////////
	//  Dependencies

	dependencies := make(map[string]*chart.Dependency)
	for _, d := range chartRequested.Metadata.Dependencies {
		dependencies[d.Alias] = d
	}

	manifests.dependencies = dependencies

	return manifests
}

func (m *Manifests) Get(kind, name string) []byte {
	step1, ok := m.manifests[kind]
	if !ok {
		return nil
	}
	step2, ok := step1[name]
	if !ok {
		return nil
	}
	return step2
}

func (m *Manifests) Set(kind, name string, manifest []byte) {
	if manifest == nil {
		log.Fatal("manifest is empty value")
	}
	if check := strings.TrimSpace(string(manifest)); check == "" {
		log.Fatal("manifest is empty value")
	}
	if m.manifests[kind] == nil {
		m.manifests[kind] = make(map[string][]byte)
	}
	m.manifests[kind][name] = manifest
}

func readValuesFile(valuesFile string) (map[string]interface{}, error) {
	bytes, err := os.ReadFile(valuesFile)
	if err != nil {
		return nil, err
	}

	vals := map[string]interface{}{}
	if err := yaml.Unmarshal(bytes, &vals); err != nil {
		return nil, err
	}

	return vals, nil
}

func SplitYAML(data []byte) [][]byte {
	var parts [][]byte
	docs := bytes.Split(data, []byte("\n---\n"))
	for _, doc := range docs {
		trimmed := bytes.TrimSpace(doc)
		if len(trimmed) > 0 {
			parts = append(parts, trimmed)
		}
	}
	return parts
}

type KubernetesManifest struct {
	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Metadata   struct {
		Name      string `yaml:"name"`
		Namespace string `yaml:"namespace"`
	} `yaml:"metadata"`
}
