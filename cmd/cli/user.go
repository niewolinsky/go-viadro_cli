package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

// userCmd represents the user command
var UserCmd = &cobra.Command{
	Use:   "user",
	Short: "Manage user's credentials",
	Long:  ``,
	Run: func(cli *cobra.Command, args []string) {
		fmt.Println("Available subcommands: register, activate, auth")
	},
}
