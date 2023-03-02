package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

// userCmd represents the user command
var AdminUserListCmd = &cobra.Command{
	Use:   "list",
	Short: "Manage users and documents as administrator",
	Long:  ``,
	Run: func(cli *cobra.Command, args []string) {
		fmt.Println("user list")
	},
}

func init() {
	AdminUserCmd.AddCommand(AdminUserListCmd)
}
