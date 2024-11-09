package main

import (
	"fmt"

	"github.com/spf13/pflag"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/resource"
	"k8s.io/kubectl/pkg/scheme"
)

func main() {
	configFlags := genericclioptions.NewConfigFlags(true)
	configFlags.AddFlags(pflag.CommandLine)
	pflag.Parse()

	builder := resource.NewBuilder(configFlags)
	err := builder.
		WithScheme(scheme.Scheme, scheme.Scheme.PrioritizedVersionsAllGroups()...).
		ResourceTypeOrNameArgs(true, pflag.Args()...).
		Flatten().
		Do().Visit(func(info *resource.Info, err error) error {
			if err != nil { panic(err) }

			fmt.Printf("namespace=%s, name=%s, obj=%T\n",
				info.Namespace, info.Name, info.Object)
			return nil
		})
	if err != nil { panic(err) }
}
