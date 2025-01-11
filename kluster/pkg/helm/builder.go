package helm

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/downloader"
	"helm.sh/helm/v3/pkg/getter"
	"sigs.k8s.io/yaml"

	"github.com/ultary/polykube/kluster/pkg/utils"
)

func Build(chartPath string, values map[string]interface{}, releaseName, namespace string) (manifests [][]byte) {
	log.Infof("[ParseFromPath] %s / %s", chartPath, namespace)

	settings := cli.New()
	actionConfig := new(action.Configuration)
	if err := actionConfig.Init(
		settings.RESTClientGetter(),
		settings.Namespace(),
		os.Getenv("HELM_DRIVER"),
		log.Debugf,
	); err != nil {
		log.Fatalf("Error initializing Helm configuration: %v\n", err)
	}

	install := action.NewInstall(actionConfig)
	install.DryRun = true
	install.ClientOnly = true
	install.Namespace = namespace
	install.ReleaseName = releaseName

	chartDownloader := downloader.Manager{
		Out:              NewLogWriter(),
		ChartPath:        chartPath,
		Keyring:          install.ChartPathOptions.Keyring,
		SkipUpdate:       false,
		Getters:          getter.All(settings),
		RepositoryConfig: settings.RepositoryConfig,
		RepositoryCache:  settings.RepositoryCache,
	}

	if err := chartDownloader.Update(); err != nil {
		log.Fatal(err)
	}

	chartRequested, err := loader.Load(chartPath)
	if err != nil {
		log.Fatal(err)
	}

	rel, err := install.Run(chartRequested, values)
	if err != nil {
		log.Fatalf("Error running install: %v\n", err)
	}

	log.Print(rel.Manifest)

	manifests = utils.SplitYAML([]byte(rel.Manifest))
	return
}

func BuildFromFileSystem(chartFS fs.FS, values map[string]interface{}, releaseName, namespace string) (manifests [][]byte) {

	tempDir, err := os.MkdirTemp("", "otlp-*")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	err = fs.WalkDir(chartFS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Fatal(err)
		}

		if !d.IsDir() {
			var data []byte
			if data, err = fs.ReadFile(chartFS, path); err != nil {
				return err
			}
			destPath := filepath.Join(tempDir, path)
			if err = os.MkdirAll(filepath.Dir(destPath), 0700); err != nil {
				return err
			}
			if err = os.WriteFile(destPath, data, 0600); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	manifests = Build(tempDir, values, releaseName, namespace)
	return
}

func BuildFromRepository() {

	const namespace = "argo"

	// TODO: this is test code for installer

	repoURL := "https://argoproj.github.io/argo-helm/"
	chartName := "argo-cd"
	chartVersion := "7.4.4"

	settings := cli.New()
	settings.SetNamespace(namespace)

	actionConfig := new(action.Configuration)
	if err := actionConfig.Init(settings.RESTClientGetter(), settings.Namespace(), "memory", log.Debugf); err != nil {
		log.Fatalf("Failed to initialize helm action: %v", err)
	}

	install := action.NewInstall(actionConfig)
	//install.ChartPathOptions.RepoURL = repoURL
	//install.ChartPathOptions.Version = chartVersion
	install.ClientOnly = true
	install.DryRun = true
	install.KubeVersion, _ = chartutil.ParseKubeVersion("v1.29.0")
	install.Namespace = settings.Namespace()
	install.ReleaseName = "argo"
	install.RepoURL = repoURL
	install.Version = chartVersion

	chartPath, err := install.ChartPathOptions.LocateChart(chartName, settings)
	if err != nil {
		log.Fatalf("Failed to locate chart: %v", err)
	}

	chart, err := loader.Load(chartPath)
	if err != nil {
		log.Fatalf("Failed to load chart: %v", err)
	}

	values := map[string]interface{}{} // Provide values for chart here if needed
	if err = yaml.Unmarshal([]byte(argoValuesYaml), &values); err != nil {
		log.Fatalf("Failed to unmarshal argo values: %v", err)
	}
	release, err := install.Run(chart, values)
	if err != nil {
		log.Fatalf("Failed to run helm install: %v", err)
	}

	fmt.Println(release.Manifest)
}

const argoValuesYaml = `
fullnameOverride: argocd
configs:
  params:
    server.insecure: true
  secret:
    createSecret: false
dex:
  enabled: false
---`
