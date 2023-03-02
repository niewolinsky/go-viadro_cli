package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

// userCmd represents the user command
var AdminUserCmd = &cobra.Command{
	Use:   "user",
	Short: "Manage users and documents as administrator",
	Long:  ``,
	Run: func(cli *cobra.Command, args []string) {
		fmt.Println("Available subcommands: delete, list, grant")
	},
}

func init() {
	AdminCmd.AddCommand(AdminUserCmd)
}
