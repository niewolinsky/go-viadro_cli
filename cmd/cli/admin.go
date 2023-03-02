package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

// userCmd represents the user command
var AdminCmd = &cobra.Command{
	Use:   "admin",
	Short: "Manage users and documents as administrator",
	Long:  ``,
	Run: func(cli *cobra.Command, args []string) {
		fmt.Println("Available subcommands: user, document")
	},
}
