package helm

import (
	"io/fs"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"github.com/ultary/monokube/kluster/pkg/utils"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/downloader"
	"helm.sh/helm/v3/pkg/getter"
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
