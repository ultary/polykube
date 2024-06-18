package helm

import (
	"fmt"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/engine"
	"sigs.k8s.io/yaml"

	"ultary.co/kluster/pkg/apps"
	"ultary.co/kluster/pkg/apps/net"
	"ultary.co/kluster/pkg/utils"
)

func Parse(chartPath string, namespace string) map[string]apps.Resource {

	releaseName := "monokube"
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

	manifests := apps.NewManifests()

	// Print the rendered templates
	for _, content := range rendered {
		blocks := utils.SplitYAML([]byte(content))
		for _, block := range blocks {
			var m KubernetesManifest
			yaml.Unmarshal(block, &m)
			manifests.Set(m.Kind, m.Metadata.Name, block)
		}
	}

	retval := map[string]apps.Resource{
		"gateway":  net.NewGateway(manifests),
		"kafka":    apps.NewKafka(manifests),
		"minio":    apps.NewMinIO(manifests),
		"postgres": apps.NewPostgreSQL(manifests),
	}

	////////////////////////////////////////////////////////////////
	//  Dependencies

	dependencies := make(map[string]*chart.Dependency)
	for _, d := range chartRequested.Metadata.Dependencies {
		dependencies[d.Alias] = d

		name := d.Name
		if d.Alias != "" {
			name = d.Alias
		}
		values, ok := vals[name].(map[string]interface{})
		if !ok {
			values = map[string]interface{}{}
		}

		retval[name] = NewChart(d, values, namespace)
	}

	return retval
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

type KubernetesManifest struct {
	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Metadata   struct {
		Name      string `yaml:"name"`
		Namespace string `yaml:"namespace"`
	} `yaml:"metadata"`
}
