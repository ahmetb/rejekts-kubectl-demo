package main

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/resource"
	kubectlget "k8s.io/kubectl/pkg/cmd/get"
)

func main() {
	configFlags := genericclioptions.NewConfigFlags(true)
	configFlags.AddFlags(pflag.CommandLine)

	// register the print flags of "kubectl get"
	printFlags := addPrintFlags(pflag.CommandLine)

	pflag.Parse()

	// initialize printer
	printer, err := printFlags.ToPrinter()
	if err != nil { panic(err) }


	builder := resource.NewBuilder(configFlags)
	err = builder.Unstructured().
		ResourceTypeOrNameArgs(true, pflag.Args()...).
		Do().Visit(func(info *resource.Info, err error) error {
			if err != nil { panic(err) }

			// Print the response object
			return printer.PrintObj(info.Object, os.Stdout)
		})
	if err != nil { panic(err) }
}


func addPrintFlags(flagSet *pflag.FlagSet) *kubectlget.PrintFlags {
	dummyCobraCmd := &cobra.Command{}
	printFlags := kubectlget.NewGetPrintFlags()
	printFlags.AddFlags(dummyCobraCmd)
	flagSet.AddFlagSet(dummyCobraCmd.Flags())
	return printFlags
}
