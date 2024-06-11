package main

import (
	"fmt"
	"os"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/engine"
	"sigs.k8s.io/yaml"
)

func main() {
	//chartPath := "path/to/chart"
	chartPath := "/Users/ghilbut/work/workbench/ultary/monokube/.helm"
	releaseName := "monokube"
	namespace := "monokube"
	valuesFile := "/Users/ghilbut/work/workbench/ultary/monokube/.helm/values.yaml"

	// Load the chart
	chart, err := loader.Load(chartPath)
	if err != nil {
		fmt.Printf("Error loading chart: %v\n", err)
		os.Exit(1)
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
	valuesToRender, err := chartutil.ToRenderValues(chart, vals, chartutil.ReleaseOptions{
		Name:      releaseName,
		Namespace: namespace,
	}, nil)
	if err != nil {
		fmt.Printf("Error creating render values: %v\n", err)
		os.Exit(1)
	}

	// Render the chart
	rendered, err := engine.Render(chart, valuesToRender)
	if err != nil {
		fmt.Printf("Error rendering chart: %v\n", err)
		os.Exit(1)
	}

	// Print the rendered templates
	for name, content := range rendered {
		fmt.Printf("--- %s ---\n%s\n", name, content)
	}

	////////////////////////////////////////////////////////////////
	//  Install

	// Create a Helm action configuration
	settings := cli.New()
	actionConfig := new(action.Configuration)
	if err := actionConfig.Init(settings.RESTClientGetter(), namespace, os.Getenv("HELM_DRIVER"), debug); err != nil {
		fmt.Printf("Error initializing action configuration: %v\n", err)
		os.Exit(1)
	}

	// Install or upgrade the Helm release
	client := action.NewInstall(actionConfig)
	client.ReleaseName = releaseName
	client.Namespace = namespace
	client.CreateNamespace = true

	release, err := client.Run(chart, vals)
	if err != nil {
		fmt.Printf("Error installing/upgrading release: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Release %s applied successfully in namespace %s\n", release.Name, release.Namespace)
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

func debug(format string, v ...interface{}) {
	fmt.Printf(format, v...)
}
