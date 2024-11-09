package main

import (
	"fmt"

	"github.com/spf13/pflag"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/resource"
)

func main() {
	configFlags := genericclioptions.NewConfigFlags(true)
	configFlags.AddFlags(pflag.CommandLine)

	// 1. Define the namespace flag
	allNamespaces := pflag.BoolP("all-namespaces", "A", false, "If present, list the requested object(s) across all namespaces.")
	pflag.Parse()

	// 2. Read the namespace preference from kubeconfig
	kubeconfigNamespace, _, err := configFlags.ToRawKubeConfigLoader().Namespace()
	if err != nil { panic(err) }

	var namespace string
	if *configFlags.Namespace != "" {
		namespace = *configFlags.Namespace // -n flag takes precedence
	} else if kubeconfigNamespace != "" {
		namespace = kubeconfigNamespace // kubeconfig namespace is used if -n flag is not set
	}


	builder := resource.NewBuilder(configFlags)
	err = builder.
		NamespaceParam(namespace).DefaultNamespace().AllNamespaces(*allNamespaces).
		Unstructured().
		ResourceTypeOrNameArgs(true, pflag.Args()...).
		Flatten().
		Do().Visit(func(info *resource.Info, err error) error {
			if err != nil { panic(err) }

			fmt.Printf("namespace=%s, name=%s\n", info.Namespace, info.Name)
			return nil
		})
	if err != nil { panic(err) }
}
