package cmd

import (
	"github.com/spf13/cobra"
	"github.com/petrbouda/golang-http-client"
	"pbouda/golang-microservices-example/yoti/client"
	"fmt"
)

var retrieveCmd = &cobra.Command{
	Use:   "retrieve",
	Short: "Retrieve a value stored under the specified ID. First-Arg = ID, Second-Arg = Encryption Key",
	Example: "retrieve my-id generated-key",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			println("Retrieve command requires 2 arguments, ID and encryption-key.")
			return
		}

		// Extract information from command-line arguments
		id := args[0]
		key := args[1]

		client := &client.EncryptClient{
			HttpClient: http_client.NewHttpClient(*debug),
		}

		value, err := client.Retrieve([]byte(id), []byte(key))
		if err != nil {
			fmt.Printf("Error occured during retrieving a value: %+v \n", err)
			return
		}

		fmt.Printf("Retrieved Value: %s\n", string(value))
	},
}

func init() {
	debug = retrieveCmd.Flags().BoolP("debug", "d", false, "Enable verbose level in HTTP Client.")

	RootCmd.AddCommand(retrieveCmd)
}
