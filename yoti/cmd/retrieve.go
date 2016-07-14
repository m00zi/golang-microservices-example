package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var retrieveCmd = &cobra.Command{
	Use:   "retrieve",
	Short: "Retrieve a value stored under the specified key.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("retrieve called")
	},
}

func init() {
	RootCmd.AddCommand(retrieveCmd)
}
