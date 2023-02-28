package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

// userCmd represents the user command
var UserCmd = &cobra.Command{
	Use:   "user",
	Short: "",
	Long:  ``,
	Run: func(cli *cobra.Command, args []string) {
		fmt.Println("Avaiable subcommands: register, activate, auth")
	},
}
