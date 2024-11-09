package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/resource"
	"k8s.io/client-go/rest"
	kubectlget "k8s.io/kubectl/pkg/cmd/get"
	"k8s.io/kubectl/pkg/scheme"
)

func main() {
	configFlags := genericclioptions.NewConfigFlags(true)
	configFlags.AddFlags(pflag.CommandLine)

	// register print flags
	printFlags := addPrintFlags(pflag.CommandLine)

	pflag.Parse()

	// initialize printer
	printer, err := printFlags.ToPrinter()
	if err != nil { panic(err) }


	builder := resource.NewBuilder(configFlags).
		ResourceTypeOrNameArgs(true, pflag.Args()...)

	// If human-readable output format, enable server-side Table printing
	isHumanOutput := *printFlags.OutputFormat == "" || *printFlags.OutputFormat == "wide"

	if isHumanOutput {
		builder.
			WithScheme(scheme.Scheme, scheme.Scheme.PrioritizedVersionsAllGroups()...).
			TransformRequests(func(r *rest.Request) {
				r.SetHeader("Accept", "application/json;as=Table;v=v1;g=meta.k8s.io,application/json")
			}).Flatten()
	} else {
		builder.Unstructured()
	}

	err = builder.Do().Visit(func(info *resource.Info, err error) error {
			if err != nil { panic(err) }

			fmt.Fprintf(os.Stderr,"Object: %T\n", info.Object)

			// Print the response object
			return printer.PrintObj(info.Object, os.Stdout)
		})
	if err != nil { panic(err) }
}


// code from https://github.com/ahmetb/kubectl-pods_on/
func addPrintFlags(flagSet *pflag.FlagSet) *kubectlget.PrintFlags {
	dummyCobraCmd := &cobra.Command{}
	printFlags := kubectlget.NewGetPrintFlags()
	printFlags.AddFlags(dummyCobraCmd)
	flagSet.AddFlagSet(dummyCobraCmd.Flags())
	return printFlags
}
