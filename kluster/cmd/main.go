package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"k8s.io/client-go/util/homedir"
)

func main() {

	appName := os.Args[0]
	rootCmd := &cobra.Command{
		Use:   appName,
		Short: "Monokube's kluster CLI application",
	}

	incluster := false
	kubeconfig, kubecontext := "", "k3s"
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = filepath.Join(home, ".kube", "config")
	}
	incluster = *rootCmd.PersistentFlags().Bool("incluster", incluster, "(optional) Use incluster authentication instead of kubeconfig and kubecontext")
	kubeconfig = *rootCmd.PersistentFlags().String("kubeconfig", kubeconfig, "(optional) absolute path to the kubeconfig file")
	kubecontext = *rootCmd.PersistentFlags().String("kubecontext", kubecontext, "(optional) The name of the kubeconfig context to use")

	rootCmd.AddCommand(NewInstallCommand(incluster, kubeconfig, kubecontext))
	rootCmd.AddCommand(NewServeCommand(incluster, kubeconfig, kubecontext))

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
