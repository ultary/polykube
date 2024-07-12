package helm

import (
	log "github.com/sirupsen/logrus"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"os"
	"sigs.k8s.io/yaml"

	"github.com/ultary/monokube/kluster/pkg/k8s"
	"github.com/ultary/monokube/kluster/pkg/utils"
)

type Chart struct {
	chart.Dependency
	values map[string]interface{}
	objs   []*unstructured.Unstructured
}

func NewChart(dependency *chart.Dependency, values map[string]interface{}, namespace string) *Chart {

	retval := &Chart{
		Dependency: *dependency,
		values:     values,
	}

	settings := cli.New()
	actionConfig := new(action.Configuration)
	if err := actionConfig.Init(settings.RESTClientGetter(), settings.Namespace(), os.Getenv("HELM_DRIVER"), log.Printf); err != nil {
		log.Fatalf("Failed to initialize Helm action configuration: %v", err)
	}

	client := action.NewInstall(actionConfig)
	client.DryRun = true
	client.ClientOnly = true
	client.Namespace = namespace
	client.ReleaseName = dependency.Name
	client.ChartPathOptions.RepoURL = dependency.Repository
	client.ChartPathOptions.Version = dependency.Version

	chartPath, err := client.ChartPathOptions.LocateChart(dependency.Name, settings)
	if err != nil {
		log.Fatalf("Failed to locate chart: %v", err)
	}

	chartRequested, err := loader.Load(chartPath)
	if err != nil {
		log.Fatalf("Failed to load chart: %v", err)
	}

	release, err := client.Run(chartRequested, values)
	if err != nil {
		log.Fatalf("Failed to render chart: %v", err)
	}

	manifests := utils.SplitYAML([]byte(release.Manifest))
	retval.objs = make([]*unstructured.Unstructured, len(manifests), len(manifests))
	for i, m := range manifests {

		var raw map[string]interface{}
		if err := yaml.Unmarshal(m, &raw); err != nil {
			log.Fatalf("Failed to unmarshal YAML: %v", err)
		}

		retval.objs[i] = &unstructured.Unstructured{Object: raw}
	}

	return retval
}

func (c *Chart) Apply(ctx k8s.Context, namespace string) error {
	for _, obj := range c.objs {
		if err := k8s.ApplyUnstructured(ctx, obj, namespace); err != nil {
			log.Fatalf("Error applying unstructed manifest in dependency: %v", err)
		}
	}
	return nil
}
