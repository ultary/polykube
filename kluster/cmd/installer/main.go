package main

import (
	"context"
	"fmt"
	"os"

	"ultary.co/kluster/pkg/k8s"
	"ultary.co/kluster/pkg/monokube"
)

func main() {

	ctx := k8s.NewContext(context.Background())
	path := os.Getenv("")
	if path == "" {
		path = "/Users/ghilbut/work/workbench/ultary/monokube/.helm"
	}
	namespace := "monokube"
	monokube.Install(ctx, path, namespace)

	os.Exit(0)

	_ = `////////////////////////////////////////////////////////////////
	//  Install by Helm

	namespace := "monokube"
	releaseName := "monokube"

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

	fmt.Printf("Release %s applied successfully in namespace %s\n", release.Name, release.Namespace)`
}

func debug(format string, v ...interface{}) {
	fmt.Printf(format, v...)
}
