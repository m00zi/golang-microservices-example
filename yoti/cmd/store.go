package cmd

import (
	"github.com/spf13/cobra"
	"pbouda/golang-microservices-example/yoti/client"
	"github.com/petrbouda/golang-http-client"
)

var storeCmd = &cobra.Command{
	Use:   "store",
	Short: "Store the given value under the specified key.",
	Run: func(cmd *cobra.Command, args []string) {
		client := &client.EncryptClient{
			HttpClient: http_client.NewHttpClient(debug),
		}

	},
}

func init() {
	RootCmd.AddCommand(storeCmd)
}
