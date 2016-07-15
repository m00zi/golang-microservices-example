package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"pbouda/golang-microservices-example/discovery"
)

var debug bool

var RootCmd = &cobra.Command{
	Use:   "client",
	Short: "Store and retrieve key/value information.",
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	// All command communicates with Service Discovery, therefore is option to set the different URL
	// than URL which is set as a default on.
	url := RootCmd.Flags().StringP("discovery-url", "u", discovery.DefaultEtcdUrl, "Service Discovery URL")

	// All command communicates with Service Discovery, therefore is option to set the different URL
	// than URL which is set as a default on.
	RootCmd.Flags().BoolVarP(&debug, "debug", "d", false, "Enable verbose level in HTTP Client.")

	// Register Service discovery client
	discovery.RegisterClient(url)
}