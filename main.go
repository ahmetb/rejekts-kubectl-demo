package main

import (
	"fmt"

	"github.com/spf13/pflag"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
)

func main() {
	configFlags := genericclioptions.NewConfigFlags(true)
	configFlags.AddFlags(pflag.CommandLine)
	pflag.Parse()

	// Print Kubernetes server version
	// Build the Kubernetes client configuration
	config, err := configFlags.ToRESTConfig()
	if err != nil { panic(err) }

	// Create the Kubernetes client
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil { panic(err) }

	versionInfo, err := clientset.Discovery().ServerVersion()
	if err != nil { panic(err) }
	fmt.Printf("Kubernetes server version: %s\n", versionInfo.GitVersion)
}
