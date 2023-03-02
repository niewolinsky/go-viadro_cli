package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

// userCmd represents the user command
var AdminUserDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Manage users and documents as administrator",
	Long:  ``,
	Run: func(cli *cobra.Command, args []string) {
		fmt.Println("user delete")
	},
}

func init() {
	AdminUserCmd.AddCommand(AdminUserDeleteCmd)
}
