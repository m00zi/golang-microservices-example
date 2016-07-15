package cmd

import (
	"github.com/spf13/cobra"
	"pbouda/golang-microservices-example/yoti/client"
	"github.com/petrbouda/golang-http-client"
	"fmt"
)

var storeCmd = &cobra.Command{
	Use:   "store",
	Short: "Store the given value under the specified id. First-Arg = ID, Second-Arg = Value",
	Example: "store my-id my-value",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			println("Store command requires 2 arguments, ID and value.")
			return
		}

		// Extract information from command-line arguments
		id := args[0]
		value := args[1]

		client := &client.EncryptClient{
			HttpClient: http_client.NewHttpClient(*debug),
		}

		aesKey, err := client.Store([]byte(id), []byte(value))
		if err != nil {
			fmt.Printf("Error occured during storing a value: %+v\n", err)
			return
		}

		fmt.Printf("Generated Ecryption Key: %s\n", string(aesKey))
	},
}

func init() {
	debug = storeCmd.Flags().BoolP("debug", "d", false, "Enable verbose level in HTTP Client.")

	RootCmd.AddCommand(storeCmd)
}
