package otlp

import (
	"crypto/md5"
	_ "embed"
	"fmt"
	log "github.com/sirupsen/logrus"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	"os"
	"sigs.k8s.io/yaml"
)

//go:embed agent.yaml
var agentConfig []byte

//go:embed collector.yaml
var collectorConfig string

//go:embed collector.yaml
var consumerConfig string

type OpenTelemetry struct {
	chartUrl     string
	chartName    string
	chartVersion string
	chartValues  map[string]interface{}
}

func NewOpenTelemetry(dependency *chart.Dependency, values map[string]interface{}) *OpenTelemetry {
	return &OpenTelemetry{
		chartUrl:     dependency.Repository,
		chartName:    dependency.Name,
		chartVersion: dependency.Version,
		chartValues:  values,
	}
}

func Apply() {

	settings := cli.New()

	actionConfig := new(action.Configuration)
	if err := actionConfig.Init(settings.RESTClientGetter(), settings.Namespace(), os.Getenv("HELM_DRIVER"), log.Printf); err != nil {
		log.Fatalf("Error initializing action configuration: %v", err)
	}

	chartURL := "https://open-telemetry.github.io/opentelemetry-helm-charts"
	chartName := "opentelemetry-collector"
	chartVersion := "0.93.0"

	client := action.NewInstall(actionConfig)
	client.ChartPathOptions.RepoURL = chartURL
	client.ChartPathOptions.Version = chartVersion

	cp, err := client.ChartPathOptions.LocateChart(chartName, settings)
	if err != nil {
		log.Fatalf("Error locating chart: %v", err)
	}

	chartRequested, err := loader.Load(cp)
	if err != nil {
		log.Fatalf("Error loading chart: %v", err)
	}

	var values map[string]interface{}
	print(string(agentConfig))
	if err := yaml.Unmarshal(agentConfig, &values); err != nil {
		log.Fatalf("Error unmarshalling YAML to StatefulSet: %v", err)
	}

	// Create a new Template client
	templateClient := action.NewInstall(actionConfig)
	templateClient.DryRun = true
	templateClient.ReleaseName = "otel-agent"
	templateClient.Namespace = "monokube"

	// Run the Template action to get the rendered manifest
	release, err := templateClient.Run(chartRequested, values)
	if err != nil {
		log.Fatalf("Error templating chart: %v", err)
	}

	fmt.Println("\n################################################################\n")
	fmt.Println(release.Manifest)
	fmt.Println("\n################################################################\n")
}

func ApplyConfigMap() {

}

func ApplyCollector() {

	checksum := md5.Sum([]byte(collectorConfig))
	print(collectorConfig)
	print(fmt.Sprintf("%x", checksum))
}

func ApplyConsumer() {

	checksum := md5.Sum([]byte(consumerConfig))
	print(consumerConfig)
	print(fmt.Sprintf("%x", checksum))
}
