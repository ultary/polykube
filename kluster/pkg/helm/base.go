package helm

import (
	"fmt"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"sigs.k8s.io/yaml"

	"github.com/ultary/monokube/kluster/pkg/apps"
	"github.com/ultary/monokube/kluster/pkg/apps/net"
	"github.com/ultary/monokube/kluster/pkg/utils"
)

func Parse(chartPath string, namespace string) map[string]apps.Resource {

	retval := make(map[string]apps.Resource)
	valuesFile := filepath.Join(chartPath, "values.yaml")

	// Load the chart
	chartRequested, err := loader.Load(chartPath)
	if err != nil {
		log.Fatalf("Error loading chart: %v\n", err)
	}

	// Load values file if provided
	vals := map[string]interface{}{}
	vals, err = readValuesFile(valuesFile)
	if err != nil {
		fmt.Printf("Error reading values file: %v\n", err)
		os.Exit(1)
	}

	////////////////////////////////////////////////////////////////
	//  Dependencies
	//
	//  NOTE: After render, dependency's chart name will be changed to alias.
	//        So, loading from remote chart is failed.
	//        If you want to prevent it, get dependency before render

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

	////////////////////////////////////////////////////////////////
	//  Render

	// Helm CLI 환경 설정
	settings := cli.New()

	// Helm action configuration 초기화
	actionConfig := new(action.Configuration)
	if err := actionConfig.Init(
		genericclioptions.NewConfigFlags(true),
		settings.Namespace(),
		os.Getenv("HELM_DRIVER"),
		func(format string, v ...interface{}) {
			fmt.Sprintf(format, v...)
		},
	); err != nil {
		log.Fatalf("Error initializing Helm configuration: %v\n", err)
	}

	// Helm 설치 작업 설정
	install := action.NewInstall(actionConfig)
	install.DryRun = true
	install.ClientOnly = true
	install.Namespace = namespace
	install.ReleaseName = "monokube"

	// Chart를 설치하고 렌더링된 매니페스트 가져오기
	rel, err := install.Run(chartRequested, vals)
	if err != nil {
		log.Fatalf("Error running install: %v\n", err)
	}

	// 렌더링된 매니페스트 출력
	fmt.Println(rel.Manifest)

	manifests := apps.NewManifests()

	blocks := utils.SplitYAML([]byte(rel.Manifest))
	for _, block := range blocks {
		var raw map[string]interface{}
		if err := yaml.Unmarshal(block, &raw); err != nil {
			log.Fatalf("Failed to unmarshal YAML: %v", err)
		}
		u := &unstructured.Unstructured{Object: raw}
		manifests.Set(u.GetKind(), u.GetName(), block)
	}

	retval["gateway"] = net.NewGateway(manifests)
	retval["kafka"] = apps.NewKafka(manifests)
	retval["minio"] = apps.NewMinIO(manifests)
	retval["postgres"] = apps.NewPostgreSQL(manifests)

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
